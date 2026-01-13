package crawler

import (
	"net/url"
	"sync"
	"time"
)

type RateLimiter struct {
	mu         sync.Mutex
	lastAccess map[string]time.Time
	delay      time.Duration
}

func NewRateLimiter(delay time.Duration) *RateLimiter {
	return &RateLimiter{
		lastAccess: make(map[string]time.Time),
		delay:      delay,
	}
}

func (r *RateLimiter) Wait(rawURL string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	host := u.Host

	r.mu.Lock()
	defer r.mu.Unlock()

	lastTime, exists := r.lastAccess[host]

	if exists {
		elapsedTime := time.Since(lastTime)

		if elapsedTime < r.delay {
			diff := r.delay - elapsedTime
			time.Sleep(diff)
		}
	}

	r.lastAccess[host] = time.Now()

}
