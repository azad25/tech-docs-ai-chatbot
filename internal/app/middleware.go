package app

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a token bucket algorithm for rate limiting
type RateLimiter struct {
	rate       float64 // tokens per second
	burst      int     // maximum burst size
	tokens     float64 // current token count
	last       time.Time
	limiterMap sync.Map // IP address -> *tokenBucket
}

type tokenBucket struct {
	tokens     float64
	last       time.Time
	rate       float64
	burst      int
	tokenMutex sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate float64, burst int) *RateLimiter {
	return &RateLimiter{
		rate:  rate,
		burst: burst,
	}
}

// getTokenBucket gets or creates a token bucket for an IP address
func (rl *RateLimiter) getTokenBucket(ip string) *tokenBucket {
	bucket, exists := rl.limiterMap.Load(ip)
	if !exists {
		bucket = &tokenBucket{
			tokens: float64(rl.burst),
			last:   time.Now(),
			rate:   rl.rate,
			burst:  rl.burst,
		}
		rl.limiterMap.Store(ip, bucket)
	}
	return bucket.(*tokenBucket)
}

// allow checks if a request should be allowed based on the rate limit
func (tb *tokenBucket) allow() bool {
	tb.tokenMutex.Lock()
	defer tb.tokenMutex.Unlock()

	now := time.Now()
	timePassed := now.Sub(tb.last).Seconds()
	tb.tokens = min(float64(tb.burst), tb.tokens+timePassed*tb.rate)
	tb.last = now

	if tb.tokens < 1 {
		return false
	}

	tb.tokens--
	return true
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// RateLimitMiddleware creates a middleware that applies rate limiting
func RateLimitMiddleware(limiter *RateLimiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				ip = forwardedFor
			}

			bucket := limiter.getTokenBucket(ip)
			if !bucket.allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Rate limit exceeded. Please try again later."))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}