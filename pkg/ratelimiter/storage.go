package ratelimiter

import (
	"sync"
	"hash/fnv"
)

type shard struct {
    mu      sync.Mutex
    buckets map[string]*Bucket
}

const shardsCount = 32

type InMemoryStore struct {
	shards []shard
}

func NewInMemoryStore() *InMemoryStore {
	store := &InMemoryStore{
		shards: make([]shard, shardsCount),
	}
	
	for i := range shardsCount {
		store.shards[i] = shard{
			buckets: make(map[string]*Bucket),
		}
	}

	return store
}

func (store *InMemoryStore) GetOrCreate(userID string) *Bucket {
	shardKey := getShardKey(userID)
	store.shards[shardKey].mu.Lock()
	defer store.shards[shardKey].mu.Unlock()

	bucket, ok := store.shards[shardKey].buckets[userID]
	if !ok {
		newBucket := NewBucket()
		store.shards[shardKey].buckets[userID] = newBucket
		return newBucket
	}

	return bucket
}

func getShardKey(userID string) uint8 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(userID))
	return uint8(h.Sum32() % 32)
}