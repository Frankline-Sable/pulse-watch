package monitor

import (
	"context"
	"net/http"
	"time"
)

type Scheduler struct {
	Client      *http.Client
	Websites    []string
	WorkerCount int
	Interval    time.Duration
	OnResults   func([]Result)
}

func (s *Scheduler) Run(ctx context.Context) error {
	if s.Interval <= 0 {
		return nil
	}
	s.runChecks(ctx)

	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if (ctx.Err()) != nil {
				return ctx.Err()
			}
			s.runChecks(ctx)

		}
	}
}
func (s *Scheduler) runChecks(ctx context.Context) {
	results := CheckAll(ctx, s.Client, s.Websites, s.WorkerCount)

	if s.OnResults != nil {
		s.OnResults(results)
	}
}
