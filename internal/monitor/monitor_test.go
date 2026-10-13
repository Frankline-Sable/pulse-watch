package monitor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckSuccess(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	client := &http.Client{
		Timeout: time.Second,
	}

	result := Check(context.Background(), client, server.URL)

	if !result.Up {
		t.Error("expected website to be UP")
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", result.StatusCode)
	}

	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}
}

func TestCheckServerError(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}),
	)
	defer server.Close()

	client := &http.Client{
		Timeout: time.Second,
	}

	result := Check(context.Background(), client, server.URL)

	if result.Up {
		t.Error("expected website to be DOWN")
	}

	if result.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", result.StatusCode)
	}
}

func TestCheckTimeout(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	client := &http.Client{
		Timeout: 50 * time.Millisecond,
	}

	result := Check(context.Background(), client, server.URL)

	if result.Up {
		t.Error("expected timed-out website to be DOWN")
	}

	if result.Error == nil {
		t.Fatal("expected a timeout error")
	}
}
