package ratelimiter

import (
	"sync"
	"testing"
)

func TestGetOrCreateBucketConcurrently(t *testing.T) {
	store := NewInMemoryStore()
	var wg sync.WaitGroup

	for range 1_000_000 {
		wg.Go(func(){
			_ = store.GetOrCreate("5")
		})
	}

	wg.Wait()
}