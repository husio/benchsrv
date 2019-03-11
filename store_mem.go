package main

import (
	"context"
	"sync"
	"time"
)

func NewMemStore() Store {
	return &memstore{}
}

// memstore is an in memory implementation of the Store. Use it for testing.
type memstore struct {
	sync.Mutex
	mem []*Benchmark
}

func (s *memstore) CreateBenchmark(ctx context.Context, content, commit string) (int64, error) {
	s.Lock()
	defer s.Unlock()

	bench := &Benchmark{
		ID:      int64(len(s.mem)) + 1,
		Created: time.Now(),
		Content: content,
		Commit:  commit,
	}
	s.mem = append(s.mem, bench)
	return bench.ID, nil
}

func (s *memstore) FindBenchmark(ctx context.Context, id int64) (*Benchmark, error) {
	s.Lock()
	defer s.Unlock()
	if id > int64(len(s.mem)) {
		return nil, ErrNotFound
	}
	bench := s.mem[id-1]
	return bench, nil
}

func (s *memstore) ListBenchmarks(ctx context.Context, olderThan time.Time, limit int) ([]*Benchmark, error) {
	s.Lock()
	defer s.Unlock()

	res := make([]*Benchmark, 0, limit)
	for _, b := range s.mem {
		if b.Created.After(olderThan) {
			continue
		}
		res = append(res, b)
		if len(res) == limit {
			break
		}
	}
	return res, nil
}
