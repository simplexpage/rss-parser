package rss_parser

import (
	"context"
	"fmt"
	"sync"
)

// worker is a function that executes jobs
func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			results <- job.execute(ctx)
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}

// WorkerPool is a struct that holds the worker pool
type WorkerPool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	Done         chan struct{}
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(wCount int) WorkerPool {
	return WorkerPool{
		workersCount: wCount,
		jobs:         make(chan Job, wCount),
		results:      make(chan Result, wCount),
		Done:         make(chan struct{}),
	}
}

// Run starts the worker pool
func (wp WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

// Results returns the results channel
func (wp WorkerPool) Results() <-chan Result {
	return wp.results
}

// AddJob adds a job to the worker pool
func (wp WorkerPool) AddJob(urls []string) {
	defer close(wp.jobs)

	for _, url := range urls {
		wp.jobs <- Job{Url: url}
	}
}
