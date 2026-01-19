package player_poller

import (
	"sync"
)

type WorkerPool[R any] struct {
	workers int
	jobs    []func() R
}

func NewWorkerPool[R any](workers int) *WorkerPool[R] {
	if workers < 1 {
		workers = 1
	}
	return &WorkerPool[R]{
		workers: workers,
		jobs:    make([]func() R, 0),
	}
}

func (p *WorkerPool[R]) Add(job func() R) {
	p.jobs = append(p.jobs, job)
}

func (p *WorkerPool[R]) ExecuteJobs() []R {
	jobCh := make(chan func() R)
	results := make(chan R, len(p.jobs))

	var wg sync.WaitGroup

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				results <- job()
			}
		}()
	}

	go func() {
		for _, job := range p.jobs {
			jobCh <- job
		}
		close(jobCh)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]R, 0, len(p.jobs))
	for r := range results {
		out = append(out, r)
	}

	return out
}
