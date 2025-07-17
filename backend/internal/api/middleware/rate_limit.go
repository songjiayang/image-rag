package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple rate limiter
// This is a basic implementation - consider using redis for production

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Clean old requests
	requests := rl.requests[ip]
	var validRequests []time.Time
	for _, req := range requests {
		if now.Sub(req) < rl.window {
			validRequests = append(validRequests, req)
		}
	}

	// Check if limit exceeded
	if len(validRequests) >= rl.limit {
		rl.requests[ip] = validRequests
		return false
	}

	// Add new request
	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests
	return true
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	rateLimiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rateLimiter.Allow(ip) {
			c.JSON(429, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": window.Seconds(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// DefaultRateLimit returns a middleware with sensible defaults
func DefaultRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(100, time.Minute)
}
