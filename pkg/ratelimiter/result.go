package ratelimiter

import "time"

type LimitResult struct {
	Allowed    bool
    Remaining  int64
    RetryAfter time.Duration
}