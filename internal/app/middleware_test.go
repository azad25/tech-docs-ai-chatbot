package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	testCases := []struct {
		name           string
		rate          float64
		burst         int
		requestCount  int
		expectedCodes []int
		waitBetween   time.Duration
	}{
		{
			name:          "allow requests within limit",
			rate:          2,
			burst:         2,
			requestCount:  2,
			expectedCodes: []int{200, 200},
			waitBetween:   0,
		},
		{
			name:          "block requests over limit",
			rate:          2,
			burst:         2,
			requestCount:  3,
			expectedCodes: []int{200, 200, 429},
			waitBetween:   0,
		},
		{
			name:          "allow requests after token replenishment",
			rate:          2,
			burst:         2,
			requestCount:  4,
			expectedCodes: []int{200, 200, 200, 200},
			waitBetween:   time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			limiter := NewRateLimiter(tc.rate, tc.burst)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := RateLimitMiddleware(limiter)
			handlerWithRateLimit := middleware(handler)

			for i := 0; i < tc.requestCount; i++ {
				if tc.waitBetween > 0 && i > 0 {
					time.Sleep(tc.waitBetween)
				}

				req := httptest.NewRequest("GET", "/", nil)
				req.RemoteAddr = "127.0.0.1:12345"
				w := httptest.NewRecorder()

				handlerWithRateLimit.ServeHTTP(w, req)
				assert.Equal(t, tc.expectedCodes[i], w.Code)
			}
		})
	}
}

func TestRateLimiterWithForwardedIP(t *testing.T) {
	limiter := NewRateLimiter(2, 2)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RateLimitMiddleware(limiter)
	handlerWithRateLimit := middleware(handler)

	// Test with X-Forwarded-For header
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	w := httptest.NewRecorder()

	// First two requests should succeed
	handlerWithRateLimit.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	handlerWithRateLimit.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Third request should be rate limited
	w = httptest.NewRecorder()
	handlerWithRateLimit.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRateLimiterConcurrent(t *testing.T) {
	limiter := NewRateLimiter(10, 10)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RateLimitMiddleware(limiter)
	handlerWithRateLimit := middleware(handler)

	// Test concurrent requests
	concurrentRequests := 20
	doneChan := make(chan bool)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = "127.0.0.1:12345"
			w := httptest.NewRecorder()

			handlerWithRateLimit.ServeHTTP(w, req)
			doneChan <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < concurrentRequests; i++ {
		<-doneChan
	}
}