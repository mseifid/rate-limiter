package ratelimiter

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

func TestLimiterAllow(t *testing.T) {
	limiter := NewLimiter(NewInMemoryStore())
	const userID string = "4"
	ctx := context.WithValue(context.Background(), "userID", userID)
	var wg sync.WaitGroup

	counter := int64(0)

	for range 1_000_000 {
		wg.Go(func() {
			res, err := limiter.Allow(ctx)
			if err != nil {
				t.Errorf("Allow function returned error")
			}
			if res.Allowed {
				atomic.AddInt64(&counter, 1)
			}
		})
	}

	wg.Wait()

	if counter < UserBucketCapacity {
		t.Errorf("limiter allow test failed, expected: %v, actual: %v", UserBucketCapacity, counter)
	}
}
