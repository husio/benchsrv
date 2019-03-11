package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	conf := configuration{
		Port:     env("PORT", "8000"),
		Secret:   env("SECRET", ""),
		Postgres: env("DATABASE_URL", ""),
	}
	if err := run(conf); err != nil {
		log.Fatal(err)
	}
}

type configuration struct {
	Port     string
	Secret   string
	Postgres string
}

func run(conf configuration) error {
	var store Store = NewMemStore()
	if conf.Postgres != "" {
		db, err := sql.Open("postgres", conf.Postgres)
		if err != nil {
			return fmt.Errorf("open PostgreSQL connection: %s", err)
		}
		if err := db.Ping(); err != nil {
			return fmt.Errorf("ping PostgreSQL: %s", err)
		}
		store, err = NewPostgresStore(db)
		if err != nil {
			return fmt.Errorf("new postgres store: %s", err)
		}
	} else {
		log.Print("using an in memory storage")
	}

	mux := http.NewServeMux()
	mux.Handle("/", listHandler(store))
	mux.Handle("/benchmarks/", showBenchmark(store))
	mux.Handle("/compare/", compareHandler(store))
	mux.Handle("/upload/", uploadHandler(store, conf.Secret))

	log.Printf("running HTTP server on port %s", conf.Port)
	if err := http.ListenAndServe(":"+conf.Port, mux); err != nil {
		return fmt.Errorf("http server: %s", err)
	}
	return nil
}

func env(name, fallback string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return fallback
}
