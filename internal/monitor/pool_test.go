package monitor

import (
	"context"
	"net/http"
	"sync"
)

func CheckAll(
	ctx context.Context,
	client *http.Client,
	websites []string,
	workerCount int) []Result {
	if workerCount < 1 {
		workerCount = 1
	}

	jobs := make(chan string)
	results := make(chan Result)

	var wg sync.WaitGroup

	// Start fixed number of workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for url := range jobs {
				result := Check(ctx, client, url)
				results <- result
			}
		}()
	}

	// Feed websites to the jobs channel
	go func() {
		defer close(jobs)

		for _, website := range websites {
			select {
			case jobs <- website:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Close results when all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	var collected []Result

	for result := range results {
		collected = append(collected, result)
	}
	return collected

}
