package ratelimiter

import (
	"sync"
	"hash/fnv"
)

type shard struct {
    mu      sync.RWMutex
    buckets map[string]*Bucket
}

const shardsCount = 32

type InMemoryStore struct {
	shards []shard
	global *Bucket
}

func NewInMemoryStore() *InMemoryStore {
	store := &InMemoryStore{
		shards: make([]shard, shardsCount),
		global: NewBucket(BucketTypeGlobal),
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
		newBucket := NewBucket(BucketTypeUser)
		store.shards[shardKey].buckets[userID] = newBucket
		return newBucket
	}

	return bucket
}

func (store *InMemoryStore) GetGlobal() *Bucket {
	return store.global
}

func getShardKey(userID string) uint8 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(userID))
	return uint8(h.Sum32() % 32)
}