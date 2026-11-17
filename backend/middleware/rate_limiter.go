package middleware

import (
	"sync"
	"time"
)

// RateLimiter limits requests per client IP using a sliding window of attempt timestamps.
type RateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	limit    int
	window   time.Duration
	getNow   func() time.Time // Allows current time to be mocked in tests
}

// Creates a new RateLimiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		getNow:   time.Now,
	}

	go limiter.purge()

	return limiter
}

// Creates a new RateLimiter. Allows injection of internal dependencies to allow for mocking.
func NewTestRateLimiter(limit int, window time.Duration, getNow func() time.Time) *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		getNow:   getNow,
	}
}

// Checks whether attempt history is stored for the given IP.
func (r *RateLimiter) HasAttempts(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.attempts[ip]
	return ok
}

// Removes history for client IPs that do not have any requests within the current window.
func (r *RateLimiter) RemoveExpired() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.getNow()
	cutoff := now.Add(-r.window)

	for ip, timestamps := range r.attempts {
		active := false

		for _, timestamp := range timestamps {
			if timestamp.After(cutoff) {
				active = true
				break
			}
		}

		if !active {
			delete(r.attempts, ip)
		}
	}
}

/*
Checks whether a request from the given IP is permitted.

The request is allowed if the IP has made less than r.limit requests in the last r.window time.
*/
func (r *RateLimiter) IsAllowed(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.getNow()
	cutoff := now.Add(-r.window)

	recent := make([]time.Time, 0, len(r.attempts[ip])+1)
	for _, timestamp := range r.attempts[ip] {
		if timestamp.After(cutoff) {
			recent = append(recent, timestamp)
		}
	}

	if len(recent) >= r.limit {
		r.attempts[ip] = recent
		return false
	}

	r.attempts[ip] = append(recent, now)
	return true
}

// Removes history for inactive client IPs every time the window expires.
func (r *RateLimiter) purge() {
	for true {
		time.Sleep(r.window)
		r.RemoveExpired()
	}
}
