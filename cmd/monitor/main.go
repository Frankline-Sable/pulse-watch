package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Frankline-Sable/pulse-watch/internal/monitor"
)

func main() {
	websites := []string{
		"https://google.com",
		"https://github.com",
		"https://example.com",
		"https://httpstat.us/503",
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	/*results := make(chan monitor.Result, len(websites))

	var wg sync.WaitGroup

	for _, website := range websites {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			result := monitor.Check(
				context.Background(),
				client,
				url,
			)

			results <- result
		}(website)
	}
	go func() {
		wg.Wait()
		close(results)
	}()*/
	results := monitor.CheckAll(
		context.Background(),
		client,
		websites, 3)

	for _, result := range results {
		if result.Error != nil {
			fmt.Printf(
				"DOWN | %s | Error: %v\n",
				result.URL,
				result.Error,
			)
			continue
		}

		status := "UP"

		if !result.Up {
			status = "DOWN"
		}
		fmt.Printf(
			"%s | %s | HTTP %d | %v\n",
			status,
			result.URL,
			result.StatusCode,
			result.ResponseTime.Round(time.Millisecond),
		)
	}
}
