package config

import "time"

type RateLimitConfig struct {
	// Login rate limiting
	LoginRateLimitKey  string
	LoginRateLimitTime time.Duration
	LoginRateLimit     int

	// Default rate limiting
	DefaultRateLimitKey  string
	DefaultRateLimitTime time.Duration
	DefaultRateLimit     int
}
