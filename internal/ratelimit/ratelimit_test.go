package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestAllowWithinBurst(t *testing.T) {
	l := New(5, 2) // burst 5, refill 2/sec
	for i := 0; i < 5; i++ {
		if !l.Allow("ip1") {
			t.Fatalf("request %d within burst should be allowed", i+1)
		}
	}
}

func TestExceedsBurst(t *testing.T) {
	l := New(5, 2)
	for i := 0; i < 5; i++ {
		l.Allow("ip1")
	}
	if l.Allow("ip1") {
		t.Fatal("request beyond burst should be denied")
	}
}

func TestIndependentKeys(t *testing.T) {
	l := New(3, 1)
	for i := 0; i < 3; i++ {
		l.Allow("a")
	}
	if l.Allow("a") {
		t.Fatal("key a should be exhausted")
	}
	if !l.Allow("b") {
		t.Fatal("key b should be independent and allowed")
	}
}

func TestRefillsOverTime(t *testing.T) {
	l := New(5, 1) // refill 1/sec
	l.now = func() time.Time { return time.Unix(0, 0) }
	for i := 0; i < 5; i++ {
		l.Allow("ip")
	}
	if l.Allow("ip") {
		t.Fatal("burst exhausted")
	}
	// advance 3 seconds: refill 3 tokens
	l.now = func() time.Time { return time.Unix(3, 0) }
	for i := 0; i < 3; i++ {
		if !l.Allow("ip") {
			t.Fatalf("refilled token %d should be allowed", i+1)
		}
	}
	if l.Allow("ip") {
		t.Fatal("only 3 refilled tokens should be available")
	}
}

func TestRefillCapsAtBurst(t *testing.T) {
	l := New(2, 5) // burst 2, refill fast
	l.now = func() time.Time { return time.Unix(0, 0) }
	l.Allow("ip")
	l.Allow("ip")
	// long idle: refill should cap at capacity (no debt beyond burst)
	l.now = func() time.Time { return time.Unix(100, 0) }
	if !l.Allow("ip") {
		t.Fatal("should be allowed after idle refill to burst")
	}
	if !l.Allow("ip") {
		t.Fatal("second after refill to burst should be allowed")
	}
	if l.Allow("ip") {
		t.Fatal("tokens must cap at burst, exceeding must be denied")
	}
}

func TestAllowZeroCapacity(t *testing.T) {
	l := New(0, 0)
	if l.Allow("ip") {
		t.Fatal("zero capacity must never allow")
	}
}

func TestConcurrentSafe(t *testing.T) {
	l := New(1000, 1000)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				l.Allow("shared")
			}
		}()
	}
	wg.Wait()
}