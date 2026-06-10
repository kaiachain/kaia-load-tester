package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/kaiachain/kaia-load-tester/klayslave/config"
	"github.com/myzhan/boomer"
)

// A StepRateLimiter is a token-bucket rate limiter whose threshold follows a
// multi-step RPS schedule (e.g. 1000 rps for 60s, then 5000 rps for 120s, ...).
// It implements the boomer.RateLimiter interface. Start() is invoked by the
// boomer runner when the locust master starts the test, so the schedule clock
// is aligned with the test start, and it restarts from the first step on every
// new test run.
type StepRateLimiter struct {
	steps        []config.RPSStep
	refillPeriod time.Duration

	currentThreshold int64
	broadcastChannel atomic.Value // holds chan bool; closed on every refill to release blocked Acquire()

	mu          sync.Mutex
	quitChannel chan bool
	doneChannel chan bool
}

// NewStepRateLimiter returns a StepRateLimiter. The steps must be non-empty
// and validated by config.ParseRPSSchedule beforehand.
func NewStepRateLimiter(steps []config.RPSStep, refillPeriod time.Duration) *StepRateLimiter {
	limiter := &StepRateLimiter{
		steps:        steps,
		refillPeriod: refillPeriod,
	}
	limiter.broadcastChannel.Store(make(chan bool))
	return limiter
}

// stepIndexAt returns the schedule step index active at the given elapsed time
// since the test started. The last step holds regardless of its duration.
func (limiter *StepRateLimiter) stepIndexAt(elapsed time.Duration) int {
	var boundary time.Duration
	for i, step := range limiter.steps {
		if i == len(limiter.steps)-1 || step.Duration == 0 {
			return i
		}
		boundary += step.Duration
		if elapsed < boundary {
			return i
		}
	}
	return len(limiter.steps) - 1
}

// Start launches the refill goroutine. It is called by the boomer runner on
// every "spawn" message from the locust master; a previous refill goroutine
// (if any) is stopped first so the schedule restarts from the first step.
func (limiter *StepRateLimiter) Start() {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	limiter.stopRefillLocked()

	quitChannel := make(chan bool)
	doneChannel := make(chan bool)
	limiter.quitChannel = quitChannel
	limiter.doneChannel = doneChannel

	go func() {
		defer close(doneChannel)
		startTime := time.Now()
		lastStep := -1
		for {
			select {
			case <-quitChannel:
				// Release any Acquire() blocked on the current bucket.
				close(limiter.broadcastChannel.Load().(chan bool))
				limiter.broadcastChannel.Store(make(chan bool))
				return
			default:
			}

			stepIdx := limiter.stepIndexAt(time.Since(startTime))
			if stepIdx != lastStep {
				step := limiter.steps[stepIdx]
				if step.Duration > 0 {
					log.Printf("StepRateLimiter: entering step %d/%d, rps=%d for %v", stepIdx+1, len(limiter.steps), step.RPS, step.Duration)
				} else {
					log.Printf("StepRateLimiter: entering step %d/%d, rps=%d until the test stops", stepIdx+1, len(limiter.steps), step.RPS)
				}
				lastStep = stepIdx
			}

			atomic.StoreInt64(&limiter.currentThreshold, limiter.steps[stepIdx].RPS)
			time.Sleep(limiter.refillPeriod)
			oldChannel := limiter.broadcastChannel.Load().(chan bool)
			limiter.broadcastChannel.Store(make(chan bool))
			close(oldChannel)
		}
	}()
}

// Acquire a token from the bucket, returns true if the bucket is exhausted.
// Exhausted callers block until the next refill.
func (limiter *StepRateLimiter) Acquire() (blocked bool) {
	permit := atomic.AddInt64(&limiter.currentThreshold, -1)
	if permit < 0 {
		blocked = true
		// block until the bucket is refilled
		<-limiter.broadcastChannel.Load().(chan bool)
	} else {
		blocked = false
	}
	return blocked
}

// Stop the rate limiter. Called by the boomer runner when the test is stopped.
func (limiter *StepRateLimiter) Stop() {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	limiter.stopRefillLocked()
}

// stopRefillLocked stops the running refill goroutine and waits for it to
// exit. The caller must hold limiter.mu.
func (limiter *StepRateLimiter) stopRefillLocked() {
	if limiter.quitChannel == nil {
		return
	}
	close(limiter.quitChannel)
	<-limiter.doneChannel
	limiter.quitChannel = nil
	limiter.doneChannel = nil
}

// runBoomerWithRateLimiter mirrors the legacy boomer.Run() but injects a
// custom rate limiter, which the legacy entrypoint cannot do. It keeps the
// Events-based stats reporting used by the testcase package working by
// forwarding "request_success"/"request_failure" events to this instance.
func runBoomerWithRateLimiter(masterHost string, masterPort int, rateLimiter boomer.RateLimiter, tasks ...*boomer.Task) {
	b := boomer.NewBoomer(masterHost, masterPort)
	b.SetRateLimiter(rateLimiter)

	boomer.Events.Subscribe("request_success", func(requestType string, name string, responseTime interface{}, responseLength int64) {
		b.RecordSuccess(requestType, name, convertResponseTime(responseTime), responseLength)
	})
	boomer.Events.Subscribe("request_failure", func(requestType string, name string, responseTime interface{}, exception string) {
		b.RecordFailure(requestType, name, convertResponseTime(responseTime), exception)
	})

	b.Run(tasks...)

	quitByMe := false
	quitChan := make(chan bool)
	boomer.Events.SubscribeOnce("boomer:quit", func() {
		if !quitByMe {
			close(quitChan)
		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-c:
		quitByMe = true
		b.Quit()
	case <-quitChan:
	}

	log.Println("shut down")
}

// convertResponseTime accepts int64 or float64 response times, matching the
// types published by the testcase package via boomer.Events.
func convertResponseTime(origin interface{}) int64 {
	switch v := origin.(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	default:
		log.Printf("unexpected responseTime type %T, recording 0", origin)
		return 0
	}
}
