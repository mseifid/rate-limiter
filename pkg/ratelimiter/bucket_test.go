package ratelimiter

import (
	"testing"
)

func TestNewBucketCreation(t *testing.T) {
	actual := NewBucket(BucketTypeUser)
	if actual.tokens != UserBucketCapacity {
		t.Errorf("user tokens are not set correctly!")
	}

	actual = NewBucket(BucketTypeGlobal)
	if actual.tokens != GlobalBucketCapacity {
		t.Errorf("global tokens are not set correctly!")
	}
}