// Package ratelimit provides a fixed-window / token-bucket style rate limiter
// keyed by an arbitrary string (e.g. client IP or device key). It is safe for
// concurrent use.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a token-bucket rate limiter per key. It is safe for concurrent use.
type Limiter struct {
	mu sync.Mutex
	// capacity is the maximum token balance (burst) per key.
	capacity float64
	// refillPerSec is tokens added per second per key.
	refillPerSec float64
	// buckets holds the current token balance per key.
	buckets map[string]float64
	// lastRefill holds the last token refill time per key.
	lastRefill map[string]time.Time
	now        func() time.Time
}

// New creates a Limiter. capacity is the burst; refillPerSec is the steady
// refill rate. capacity=0 denies everything; refillPerSec<=0 gives burst only.
func New(capacity, refillPerSec float64) *Limiter {
	return &Limiter{
		capacity:    capacity,
		refillPerSec: refillPerSec,
		buckets:     make(map[string]float64),
		lastRefill:  make(map[string]time.Time),
		now:         time.Now,
	}
}

// Allow reports whether key may proceed, consuming one token when true.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.capacity <= 0 {
		return false
	}

	now := l.now()
	tokens, ok := l.buckets[key]
	if !ok {
		tokens = l.capacity
		l.buckets[key] = tokens
	}
	last := l.lastRefill[key]
	if !last.IsZero() {
		elapsed := now.Sub(last).Seconds()
		if elapsed > 0 {
			add := elapsed * l.refillPerSec
			tokens += add
			if tokens > l.capacity {
				tokens = l.capacity
			}
		}
	}
	l.lastRefill[key] = now
	l.buckets[key] = tokens

	if tokens >= 1 {
		l.buckets[key] = tokens - 1
		return true
	}
	return false
}