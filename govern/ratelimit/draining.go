package ratelimit

import ()

func NewRateLimitOnDraining(rate int) *RateLimitOnBucket {
	return NewRateLimitOnBucket(1, rate)
}
