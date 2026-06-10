package main

import (
	"testing"
	"time"

	"github.com/kaiachain/kaia-load-tester/klayslave/config"
)

func TestStepRateLimiterStepIndexAt(t *testing.T) {
	limiter := NewStepRateLimiter([]config.RPSStep{
		{RPS: 1000, Duration: 60 * time.Second},
		{RPS: 5000, Duration: 120 * time.Second},
		{RPS: 10000, Duration: 0}, // hold until stop
	}, time.Second)

	tests := []struct {
		elapsed  time.Duration
		expected int
	}{
		{0, 0},
		{59 * time.Second, 0},
		{60 * time.Second, 1},
		{179 * time.Second, 1},
		{180 * time.Second, 2},
		{24 * time.Hour, 2}, // last step holds forever
	}
	for _, tt := range tests {
		if got := limiter.stepIndexAt(tt.elapsed); got != tt.expected {
			t.Errorf("stepIndexAt(%v): expected step %d, got %d", tt.elapsed, tt.expected, got)
		}
	}

	// A single-step schedule always stays on step 0.
	single := NewStepRateLimiter([]config.RPSStep{{RPS: 100, Duration: 10 * time.Second}}, time.Second)
	if got := single.stepIndexAt(time.Hour); got != 0 {
		t.Errorf("single-step stepIndexAt(1h): expected 0, got %d", got)
	}
}

func TestStepRateLimiterAcquireAndRestart(t *testing.T) {
	limiter := NewStepRateLimiter([]config.RPSStep{
		{RPS: 3, Duration: 0},
	}, 50*time.Millisecond)

	limiter.Start()
	defer limiter.Stop()

	// Wait for the refill goroutine to publish the first threshold.
	deadline := time.Now().Add(time.Second)
	for limiter.Acquire() {
		if time.Now().After(deadline) {
			t.Fatal("limiter never allowed an acquire")
		}
	}
	// One token consumed above; the remaining two must pass without blocking.
	for i := 0; i < 2; i++ {
		if blocked := limiter.Acquire(); blocked {
			t.Fatalf("acquire %d unexpectedly blocked before the bucket was exhausted", i)
		}
	}
	// The 4th acquire must block until the next refill, then succeed.
	if blocked := limiter.Acquire(); !blocked {
		t.Fatal("acquire expected to block after the bucket was exhausted")
	}

	// Start() again (new test run from the locust UI) must restart cleanly.
	limiter.Start()
	deadline = time.Now().Add(time.Second)
	for limiter.Acquire() {
		if time.Now().After(deadline) {
			t.Fatal("limiter never allowed an acquire after restart")
		}
	}
}
