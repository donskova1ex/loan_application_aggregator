package ratelimit

import (
	"sync"
	"time"
)

type Limiter interface {
	IsAllowed(key string) bool
	GetRemaining(key string) int
	Reset(key string)
	GetStats() map[string]interface{}
}

type SlidingWindowLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	config   *Config
	stats    *Stats
}

type Stats struct {
	TotalRequests   int
	BlockedRequests int
	TopKeys         map[string]int
	mu              sync.RWMutex
}

func NewSlidingWindowLimiter(config *Config) *SlidingWindowLimiter {
	limiter := &SlidingWindowLimiter{
		requests: make(map[string][]time.Time),
		config:   config,
		stats: &Stats{
			TopKeys: make(map[string]int),
		},
	}

	go limiter.startCleanup()

	return limiter
}

func (l *SlidingWindowLimiter) IsAllowed(key string) bool {

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	threshold := now.Add(-l.config.WindowSize)
	requests, exists := l.requests[key]
	if !exists {
		requests = []time.Time{}
	}

	var validRequests []time.Time
	for _, requestTime := range requests {
		if requestTime.After(threshold) {
			validRequests = append(validRequests, requestTime)
		}
	}

	if len(validRequests) >= l.config.RequestsPerMinute {
		l.stats.recordRequest(key, true)
		return false
	}

	validRequests = append(validRequests, now)
	l.requests[key] = validRequests
	l.stats.recordRequest(key, false)
	return true
}

func (l *SlidingWindowLimiter) GetRemaining(key string) int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	now := time.Now()
	threshold := now.Add(-l.config.WindowSize)
	requests, exists := l.requests[key]
	if !exists {
		return l.config.RequestsPerMinute
	}

	var validRequests []time.Time
	for _, requestTime := range requests {
		if requestTime.After(threshold) {
			validRequests = append(validRequests, requestTime)
		}
	}

	return l.config.RequestsPerMinute - len(validRequests)
}

func (l *SlidingWindowLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.requests, key)
}

func (l *SlidingWindowLimiter) GetStats() map[string]interface{} {
	return l.stats.getStats()
}

func (l *SlidingWindowLimiter) startCleanup() {
	ticker := time.NewTicker(l.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		l.cleanup()
	}
}

func (l *SlidingWindowLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	threshold := now.Add(-l.config.WindowSize)
	for key, requests := range l.requests {
		var validRequests []time.Time
		for _, requestTime := range requests {
			if requestTime.After(threshold) {
				validRequests = append(validRequests, requestTime)
			}
		}
		switch {
		case len(validRequests) == 0:
			delete(l.requests, key)
		default:
			l.requests[key] = validRequests
		}
	}
}

func (s *Stats) recordRequest(key string, blocked bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TotalRequests++
	if blocked {
		s.BlockedRequests++
	}
	s.TopKeys[key]++
}

func (s *Stats) getStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	topKeys := make(map[string]int)
	for key, count := range s.TopKeys {
		topKeys[key] = count
	}

	return map[string]interface{}{
		"total_requests":              s.TotalRequests,
		"blocked_requests":            s.BlockedRequests,
		"blocked_requests_percentage": float64(s.BlockedRequests) / float64(s.TotalRequests) * 100,
		"top_keys":                    topKeys,
	}
}
