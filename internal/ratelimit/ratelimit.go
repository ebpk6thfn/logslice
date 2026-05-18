// Package ratelimit provides a simple token-bucket rate limiter for
// controlling the throughput of log line processing pipelines.
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Limiter controls the rate at which operations are allowed to proceed.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to ratePerSec operations per second
// with a burst capacity equal to ratePerSec (one second of tokens).
func New(ratePerSec float64) *Limiter {
	if ratePerSec <= 0 {
		ratePerSec = 1
	}
	now := time.Now()
	return &Limiter{
		tokens:   ratePerSec,
		max:      ratePerSec,
		rate:     ratePerSec,
		lastTick: now,
		clock:    time.Now,
	}
}

// Wait blocks until a token is available or ctx is cancelled.
// Returns ctx.Err() if the context is cancelled before a token is acquired.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		l.mu.Lock()
		l.refill()
		if l.tokens >= 1.0 {
			l.tokens -= 1.0
			l.mu.Unlock()
			return nil
		}
		// Calculate how long until the next token is available.
		wait := time.Duration(float64(time.Second) / l.rate)
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}

// TryAcquire attempts to acquire a token without blocking.
// Returns true if a token was acquired, false otherwise.
func (l *Limiter) TryAcquire() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill()
	if l.tokens >= 1.0 {
		l.tokens -= 1.0
		return true
	}
	return false
}

// refill adds tokens based on elapsed time. Must be called with l.mu held.
func (l *Limiter) refill() {
	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}
	l.lastTick = now
}
