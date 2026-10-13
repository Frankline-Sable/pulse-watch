package monitor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestCheckAll(t *testing.T) {
	var requests atomic.Int32

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests.Add(1)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	websites := []string{
		server.URL + "/one",
		server.URL + "/two",
		server.URL + "/three",
		server.URL + "/four",
	}

	client := &http.Client{
		Timeout: time.Second,
	}

	results := CheckAll(
		context.Background(),
		client,
		websites,
		2,
	)

	if len(results) != len(websites) {
		t.Fatalf(
			"expected %d results, got %d",
			len(websites),
			len(results),
		)
	}

	for _, result := range results {
		if !result.Up {
			t.Errorf("expected %s to be UP", result.URL)
		}
	}

	if requests.Load() != int32(len(websites)) {
		t.Errorf(
			"expected %d requests, got %d",
			len(websites),
			requests.Load(),
		)
	}
}

func TestWorkerPoolConcurrencyLimit(t *testing.T) {
	var active atomic.Int32
	var maxActive atomic.Int32

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			current := active.Add(1)
			defer active.Add(-1)

			for {
				previous := maxActive.Load()
				if current <= previous ||
					maxActive.CompareAndSwap(previous, current) {
					break
				}
			}

			time.Sleep(30 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	var websites []string

	for i := 0; i < 10; i++ {
		websites = append(websites, server.URL)
	}

	results := CheckAll(
		context.Background(),
		&http.Client{Timeout: time.Second},
		websites,
		3,
	)

	if len(results) != 10 {
		t.Fatalf("expected 10 results, got %d", len(results))
	}

	if maxActive.Load() > 3 {
		t.Errorf(
			"concurrency exceeded limit: %d",
			maxActive.Load(),
		)
	}
}
