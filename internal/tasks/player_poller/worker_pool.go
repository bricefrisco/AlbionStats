package player_poller

import (
	"sync"
	"time"
)

type Job[R any] interface {
	Run() R
}

type WorkerPool[R any] struct {
	workers int
	limiter <-chan time.Time
	jobs    []Job[R]
}

func NewWorkerPool[R any](ratePerSec int) *WorkerPool[R] {
	var limiter <-chan time.Time
	if ratePerSec > 0 {
		t := time.NewTicker(time.Second / time.Duration(ratePerSec))
		limiter = t.C
	}

	workers := ratePerSec
	if workers < 1 {
		workers = 1
	}

	return &WorkerPool[R]{
		workers: workers,
		limiter: limiter,
		jobs:    make([]Job[R], 0),
	}
}

func (p *WorkerPool[R]) Add(job Job[R]) {
	p.jobs = append(p.jobs, job)
}

func (p *WorkerPool[R]) ExecuteJobs() []R {
	jobCh := make(chan Job[R])
	results := make(chan R, len(p.jobs))

	var wg sync.WaitGroup

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				if p.limiter != nil {
					<-p.limiter
				}
				results <- job.Run()
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
