package ratelimiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewBucketCreation(t *testing.T) {
	tests := []struct {
		name           string
		bucketType     BucketType
		expectedTokens int64
	}{
		{"user bucket creation", BucketTypeUser, UserBucketCapacity},
		{"global bucket creation", BucketTypeGlobal, GlobalBucketCapacity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewBucket(tt.bucketType)
			if actual.tokens != tt.expectedTokens {
				t.Errorf("expected: %v, got: %v", tt.expectedTokens, actual.tokens)
			}
		})
	}
}

func TestConsumeBucket(t *testing.T) {
	buck := NewBucket(BucketTypeUser)
	actualRes := buck.Consume()

	if actualRes.Remaining != buck.capacity-1 {
		t.Errorf("token count did not consumed correctly, expected: %v, got: %v", buck.capacity-1, actualRes.Remaining)
	}

	if !actualRes.Allowed {
		t.Errorf("request allowance is not correct, expected: %v, got: %v", true, false)
	}
}

func TestRefillBucket(t *testing.T) {
	buck := NewBucket(BucketTypeUser)

	tests := []struct {
		name                string
		sleepDuration       time.Duration
		initToken           int64
		expectedAfterRefill int64
	}{
		{"low usage, more than enough time to be full",
			buck.oneTokenRefillDuration * 5,
			buck.capacity - 1,
			buck.capacity,
		},
		{"high usage, less than enough time to be full",
			buck.oneTokenRefillDuration * 3,
			buck.capacity - 10,
			(buck.capacity - 10) + (int64(buck.oneTokenRefillDuration * 3) / int64(buck.oneTokenRefillDuration)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buck.tokens = tt.initToken
			time.Sleep(tt.sleepDuration)
			buck.refill()
			if buck.tokens != tt.expectedAfterRefill {
				t.Errorf("expected: %v, got: %v", tt.expectedAfterRefill, buck.tokens)
			}
		})
	}
}

func TestConsumeBucketConcurrently(t *testing.T) {
	buck := NewBucket(BucketTypeUser)
	buck.oneTokenRefillDuration = time.Minute * 1 // to deactivate refilling during test
	var wg sync.WaitGroup
	allowed := int64(0)
	for range 1_000_000 {
		wg.Go(func() {
			res := buck.Consume()
			if res.Allowed {
				atomic.AddInt64(&allowed, 1)
			}
		})
	}

	wg.Wait()
	if allowed != buck.capacity {
		t.Errorf("consume token concurrently. expected: %v, actual: %v", buck.capacity, allowed)
	}
}
