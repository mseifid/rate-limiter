package ratelimiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewBucketCreation(t *testing.T) {
	actual := NewBucket(BucketTypeUser)
	if actual.tokens != actual.capacity {
		t.Errorf("user tokens are not set correctly!")
	}

	actual = NewBucket(BucketTypeGlobal)
	if actual.tokens != actual.capacity {
		t.Errorf("global tokens are not set correctly!")
	}
}

func TestConsumeBucket(t *testing.T) {
	buck := NewBucket(BucketTypeUser)
	actualRes := buck.Consume()

	if actualRes.Remaining != buck.capacity - 1 {
		t.Errorf("token count did not consumed correctly")
	}

	if !actualRes.Allowed {
		t.Errorf("request allowance is not correct")
	}
}

func TestRefillBucket(t *testing.T) {
	buck := NewBucket(BucketTypeUser)
	buck.tokens = 199
	time.Sleep(time.Second * 1)
	buck.refill()

	if buck.tokens != buck.capacity {
		t.Errorf("refill (maybe MinInt) is not working correctly")
	}

	buck.tokens = 190
	time.Sleep(time.Second * 1)
	buck.refill()
	if buck.tokens != (190 + (int64(time.Second * 1) / int64(buck.oneTokenRefillDuraion))) {
		t.Errorf("refill is not working correctly")
	}
}

func TestConsumeBucketConcurrently(t *testing.T) {
	buck := NewBucket(BucketTypeUser)
	buck.oneTokenRefillDuraion = time.Minute *1 // to deactivate refilling during test
	var wg sync.WaitGroup
	allowed := int64(0)
	for range 1_000_000 {
		wg.Go(func(){
			res := buck.Consume()
			if res.Allowed {
				atomic.AddInt64(&allowed, 1)
			}
		})
	}

	wg.Wait()
	if allowed != buck.capacity {
		t.Errorf("oh no! expected: %v, actual: %v", buck.capacity, allowed)
	}
}