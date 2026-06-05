package ratelimit

import (
	"sync"
	"time"
)

type limiter struct {
	last time.Time
}

type RateLimiter struct {
	mu sync.Mutex
	m  map[string]*limiter
}

func New() *RateLimiter {
	return &RateLimiter{
		m: make(map[string]*limiter),
	}
}

func (r *RateLimiter) Allow(ip string, interval time.Duration) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	if l, ok := r.m[ip]; ok {
		if now.Sub(l.last) < interval {
			return false
		}
		l.last = now
		return true
	}

	r.m[ip] = &limiter{last: now}
	return true
}
