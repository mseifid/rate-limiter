package ratelimiter

import "time"

type BucketType string

const (
	BucketTypeGlobal BucketType = "global"
	BucketTypeUser BucketType = "user"
)

const UserBucketCapacity = 200 // per minute
const GlobalBucketCapacity = 1000 // per minute
const UserRefillDuration = time.Millisecond * 330
const GlobalRefillDuration = time.Millisecond * 60