package monitor

import (
	"context"
	"net/http"
	"time"
)

type Result struct {
	URL          string
	StatusCode   int
	ResponseTime time.Duration
	Up           bool
	Error        error
	CheckedAt    time.Time
}

func check(ctx context.Context, client *http.Client, url string) Result {
	start := time.Now()

	result := Result{
		URL:       url,
		CheckedAt: start,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		result.Error = err
		return result
	}
	resp, err := client.Do(req)
	result.ResponseTime = time.Since(start)

	if err != nil {
		result.Error = err
		return result
	}

	defer func() {
		// We don't want to overwrite the result error
		if closeErr := resp.Body.Close(); closeErr != nil && result.Error == nil {
			result.Error = closeErr
		}
	}()

	result.StatusCode = resp.StatusCode
	result.Up = resp.StatusCode >= 200 && resp.StatusCode < 400

	return result
}
