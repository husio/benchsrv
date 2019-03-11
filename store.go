package main

import (
	"context"
	"errors"
	"time"
)

type Store interface {
	// CreateBenchmark creates and stores a new instance of a benchmark.
	// Newly create benchmark ID is returned on success.
	CreateBenchmark(ctx context.Context, content, commit string) (int64, error)

	// FindBenchmark returns a benchmark instance with given ID. If not
	// found ErrNotFound is returned instead.
	FindBenchmark(ctx context.Context, benchmarkID int64) (*Benchmark, error)

	// ListBenchmarks returns a list of benchmarks that match given
	// criteria. If no benchmarks were found an empty set is returned and
	// no error.
	ListBenchmarks(ctx context.Context, olderThan time.Time, limit int) ([]*Benchmark, error)
}

// Benchmark represents a single benchmark instance as stored in the store.
type Benchmark struct {
	// An ID of a benchmark instance.
	ID int64

	// Created represents the time that the benchmark was created
	// (uploaded) in the database. This is server creation time.
	Created time.Time

	// Content represents the output result of a single benchmark run.
	Content string

	// Commit ID that this benchmark was runnig for.
	Commit string
}

var ErrNotFound = errors.New("not found")
