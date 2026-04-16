package ratelimiter

import (
	"app/pkg/utility"
	"sync"
	"time"
)

type Bucket struct {
	mu                    sync.Mutex
	tokens                int64
	lastRefill            time.Time
	oneTokenRefillDuraion time.Duration
	capacity              int64
}

func (bucket *Bucket) Consume() LimitResult {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	if bucket.tokens == 0 && time.Since(bucket.lastRefill) < bucket.oneTokenRefillDuraion {
		return LimitResult{
			Allowed:    false,
			Remaining:  0,
			RetryAfter: bucket.oneTokenRefillDuraion - time.Since(bucket.lastRefill),
		}
	}

	bucket.refill()

	bucket.tokens--

	res := LimitResult{
		Allowed:    true,
		Remaining:  bucket.tokens,
		RetryAfter: 0,
	}

	if res.Remaining == 0 {
		res.RetryAfter = bucket.oneTokenRefillDuraion
	}

	return res
}

func NewBucket(bucketType BucketType) *Bucket {
	if bucketType == BucketTypeUser {
		return &Bucket{
			tokens:                UserBucketCapacity,
			lastRefill:            time.Now(),
			oneTokenRefillDuraion: UserRefillDuration,
			capacity:              UserBucketCapacity,
		}
	}

	return &Bucket{
		tokens:                GlobalBucketCapacity,
		lastRefill:            time.Now(),
		oneTokenRefillDuraion: GlobalRefillDuration,
		capacity:              GlobalBucketCapacity,
	}

}

func (bucket *Bucket) refill() {
	elapsedFromRefill := time.Since(bucket.lastRefill)
	tokensToAdd := int64(elapsedFromRefill / bucket.oneTokenRefillDuraion)

	bucket.tokens = utility.MinInt(bucket.capacity, bucket.tokens+tokensToAdd)

	bucket.lastRefill = bucket.lastRefill.Add(time.Duration(tokensToAdd) * bucket.oneTokenRefillDuraion)
}
