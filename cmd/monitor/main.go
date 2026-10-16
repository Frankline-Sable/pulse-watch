package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Frankline-Sable/pulse-watch/internal/monitor"
)

func main() {
	interval := flag.Duration(
		"interval",
		30*time.Second,
		"Time between monitoring checks",
	)

	once := flag.Bool(
		"once",
		false,
		"Run monitoring checks once and exit",
	)

	workers := flag.Int(
		"workers",
		3,
		"Maximum concurrent monitoring workers",
	)

	flag.Parse()

	if *workers < 1 {
		fmt.Fprintln(os.Stderr, "workers must be at least 1")
		os.Exit(1)
	}

	if !*once && *interval <= 0 {
		fmt.Fprintln(os.Stderr, "interval must be positive")
		os.Exit(1)
	}

	websites := []string{
		"https://google.com",
		"https://github.com",
		"https://example.com",
		"https://httpstat.us/503",
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Cancel the context on Ctrl+C or SIGTERM.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	printResults := func(results []monitor.Result) {
		fmt.Printf(
			"\n--- Monitoring check: %s ---\n",
			time.Now().Format(time.RFC3339),
		)

		for _, result := range results {
			status := "UP"
			if !result.Up {
				status = "DOWN"
			}

			if result.Error != nil {
				fmt.Printf(
					"%s | %s | Error: %v\n",
					status,
					result.URL,
					result.Error,
				)
				continue
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

	if *once {
		results := monitor.CheckAll(
			ctx,
			client,
			websites,
			*workers,
		)

		printResults(results)
		return
	}

	scheduler := monitor.Scheduler{
		Client:      client,
		Websites:    websites,
		WorkerCount: *workers,
		Interval:    *interval,
		OnResults:   printResults,
	}

	fmt.Printf(
		"PulseWatch started | Interval: %s | Workers: %d\n",
		*interval,
		*workers,
	)

	err := scheduler.Run(ctx)

	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, "monitoring error:", err)
		os.Exit(1)
	}

	fmt.Println("\nPulseWatch stopped gracefully.")
}
