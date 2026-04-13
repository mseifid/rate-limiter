package ratelimiter

import (
	"app/pkg/utility"
	"sync"
	"time"
)

type Bucket struct {
	mu         sync.Mutex
	tokens     int64
	lastRefill time.Time
	oneTokenRefillDuraion time.Duration
}

var globalBucket *Bucket = &Bucket{
	tokens:     maxGlobalTokens,
	lastRefill: time.Now(),
	oneTokenRefillDuraion: time.Millisecond * 60,
}

const maxUserTokens int64 = 200
const maxGlobalTokens int64 = 1000

func (bucket *Bucket) Consume() LimitResult {
	bucket.mu.Lock()
	globalBucket.mu.Lock()
	defer bucket.mu.Unlock()
	defer globalBucket.mu.Unlock()

	if bucket.tokens == 0 && time.Since(bucket.lastRefill) < bucket.oneTokenRefillDuraion {
		return LimitResult{
			Allowed:    false,
			Remaining:  0,
			RetryAfter: bucket.oneTokenRefillDuraion - time.Since(bucket.lastRefill),
		}
	}

	if globalBucket.tokens == 0 && time.Since(globalBucket.lastRefill) < globalBucket.oneTokenRefillDuraion {
		return LimitResult{
			Allowed:    false,
			Remaining:  0,
			RetryAfter: globalBucket.oneTokenRefillDuraion - time.Since(globalBucket.lastRefill),
		}
	}

	bucket.refill(maxUserTokens)
	globalBucket.refill(maxGlobalTokens)

	bucket.tokens--
	globalBucket.tokens--

	res := LimitResult{
		Allowed: true,
		Remaining: bucket.tokens,
		RetryAfter: 0,
	}

	if res.Remaining == 0 {
		res.RetryAfter = bucket.oneTokenRefillDuraion
	}

	return res
}

func NewBucket() *Bucket {
	return &Bucket{
		tokens:     maxUserTokens,
		lastRefill: time.Now(),
		oneTokenRefillDuraion: time.Millisecond * 330,
	}
}

func (bucket *Bucket) refill(maxTokensAllowed int64) {
	elapsedFromRefill := time.Since(bucket.lastRefill)
	tokensToAdd := int64(elapsedFromRefill / bucket.oneTokenRefillDuraion)

	bucket.tokens = utility.MinInt(maxTokensAllowed, bucket.tokens+tokensToAdd)

	bucket.lastRefill = bucket.lastRefill.Add(time.Duration(tokensToAdd) * bucket.oneTokenRefillDuraion)
}
