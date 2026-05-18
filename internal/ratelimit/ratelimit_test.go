package ratelimit

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNew_DefaultsToPositiveRate(t *testing.T) {
	l := New(0)
	if l.rate <= 0 {
		t.Fatalf("expected positive rate, got %v", l.rate)
	}
}

func TestTryAcquire_ConsumesToken(t *testing.T) {
	l := New(10)
	if !l.TryAcquire() {
		t.Fatal("expected first TryAcquire to succeed")
	}
}

func TestTryAcquire_ExhaustsTokens(t *testing.T) {
	// Rate of 2 means burst of 2 tokens available immediately.
	l := New(2)
	if !l.TryAcquire() {
		t.Fatal("first acquire should succeed")
	}
	if !l.TryAcquire() {
		t.Fatal("second acquire should succeed")
	}
	if l.TryAcquire() {
		t.Fatal("third acquire should fail — no tokens left")
	}
}

func TestTryAcquire_RefillsOverTime(t *testing.T) {
	l := New(100)
	// Drain all tokens.
	for l.TryAcquire() {
	}
	// Advance the internal clock by 1 second.
	l.mu.Lock()
	l.lastTick = time.Now().Add(-1 * time.Second)
	l.mu.Unlock()

	if !l.TryAcquire() {
		t.Fatal("expected token after refill")
	}
}

func TestWait_AcquiresToken(t *testing.T) {
	l := New(1000)
	ctx := context.Background()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWait_RespectsContextCancellation(t *testing.T) {
	// Rate of 0.001 means ~1000s between tokens; drain burst first.
	l := New(0.001)
	l.TryAcquire() // drain the single burst token

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if err := l.Wait(ctx); err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestWait_ConcurrentSafe(t *testing.T) {
	l := New(1000)
	var wg sync.WaitGroup
	ctx := context.Background()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Wait(ctx)
		}()
	}
	wg.Wait()
}
