package ratelimiter

import "time"

type bucketType string

const (
	bucketTypeGlobal bucketType = "global"
	bucketTypeUser bucketType = "user"
	userBucketCapacity = 200  // per minute
	globalBucketCapacity = 1000 // per minute
	userRefillDuration = time.Millisecond * 330
	globalRefillDuration = time.Millisecond * 60
)
