package stunserver

import (
	"sync"
	"time"
)

type IPRateLimiter struct {
	ips   map[string]*client
	mu    sync.Mutex
	rate  int
	burst int
}

type client struct {
	tokens     int
	lastRefill time.Time
}

func NewIPRateLimiter(rate int, interval time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		ips:   make(map[string]*client),
		rate:  rate,
		burst: rate,
	}
}

func (l *IPRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	c, exists := l.ips[ip]
	if !exists {
		l.ips[ip] = &client{
			tokens:     l.burst - 1,
			lastRefill: time.Now(),
		}
		return true
	}

	now := time.Now()
	elapsed := now.Sub(c.lastRefill).Seconds()

	refill := int(elapsed * float64(l.rate))
	if refill > 0 {
		c.tokens += refill
		if c.tokens > l.burst {
			c.tokens = l.burst
		}
		c.lastRefill = now
	}

	if c.tokens > 0 {
		c.tokens--
		return true
	}

	return false
}
