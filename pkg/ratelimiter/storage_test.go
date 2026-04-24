package ratelimiter

import (
	"sync"
	"testing"
)

// This test should be ran using -race flag to show possible race conditions
func TestGetOrCreateBucketConcurrently(t *testing.T) {
	store := NewInMemoryStore()
	var wg sync.WaitGroup
	const userID string = "5";

	for range 1_000_000 {
		wg.Add(1)
		go func(){
			defer wg.Done()
			_ = store.GetOrCreate(userID)
		}()
	}

	wg.Wait()
}