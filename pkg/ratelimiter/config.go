package ratelimiter

import "time"

type BucketType string

const (
	bucketTypeGlobal BucketType = "global"
	bucketTypeUser BucketType = "user"
	userBucketCapacity = 200  // per minute
	globalBucketCapacity = 1000 // per minute
	userRefillDuration = time.Millisecond * 330
	globalRefillDuration = time.Millisecond * 60
)
