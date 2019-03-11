package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type postgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (Store, error) {
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("cannot migrate: %s", err)
	}
	return &postgresStore{db: db}, nil
}

func (pg *postgresStore) CreateBenchmark(ctx context.Context, content, commit string) (int64, error) {
	res := pg.db.QueryRowContext(ctx, `
		INSERT INTO benchmarks (created, content, commit)
		VALUES ($1, $2, $3)
		RETURNING id
	`, time.Now(), content, commit)

	var id int64
	err := res.Scan(&id)
	return id, err
}

func (pg *postgresStore) FindBenchmark(ctx context.Context, benchID int64) (*Benchmark, error) {
	res := pg.db.QueryRowContext(ctx, `
		SELECT created, content, commit FROM benchmarks
		WHERE id = $1 LIMIT 1
	`, benchID)

	var b Benchmark
	switch err := res.Scan(&b.Created, &b.Content, &b.Commit); err {
	case sql.ErrNoRows:
		return nil, ErrNotFound
	case nil:
		return &b, nil
	default:
		return nil, err
	}
}

func (pg *postgresStore) ListBenchmarks(ctx context.Context, olderThan time.Time, limit int) ([]*Benchmark, error) {
	rows, err := pg.db.QueryContext(ctx, `
		SELECT id, created, content, commit FROM benchmarks
		WHERE created < $1
		ORDER BY created DESC
		LIMIT $2
	`, olderThan, limit)
	if err != nil {
		return nil, fmt.Errorf("cannot query benchmarks: %s", err)
	}
	defer rows.Close()

	results := make([]*Benchmark, 0, limit)
	for rows.Next() {
		var b Benchmark
		if err := rows.Scan(&b.ID, &b.Created, &b.Content, &b.Commit); err != nil {
			return nil, fmt.Errorf("cannot scan result: %s", err)
		}
		results = append(results, &b)
	}

	return results, rows.Err()
}

func migrate(db *sql.DB) error {
	for _, query := range strings.Split(schema, "\n---\n") {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("%s: %s", err, query)
		}
	}
	return nil
}

const schema = `
CREATE TABLE IF NOT EXISTS benchmarks (
	id SERIAL,
	created TIMESTAMPTZ NOT NULL,
	content TEXT NOT NULL
);

---

ALTER TABLE benchmarks ADD COLUMN IF NOT EXISTS commit TEXT NOT NULL;
`
