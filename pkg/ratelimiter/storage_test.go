package ratelimiter

import (
	"sync"
	"testing"
)

func TestGetOrCreateBucketConcurrently(t *testing.T) {
	store := NewInMemoryStore()
	var wg sync.WaitGroup
	const userID string = "5";

	for range 1_000_000 {
		wg.Go(func(){
			_ = store.GetOrCreate(userID)
		})
	}

	wg.Wait()
}