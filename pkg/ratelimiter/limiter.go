package ratelimiter

import (
	"context"
	"fmt"
)

type Limiter struct {
	store Store
}

type Store interface {
	GetOrCreate(userID string) *Bucket
}

type TokenBucket interface {
	Consume() LimitResult
}

func NewLimiter(store Store) *Limiter {
	return &Limiter{
		store: store,
	}
}

func (limiter *Limiter) Allow(ctx context.Context) (LimitResult, error) {
	userID := ctx.Value("userID").(string)
	fmt.Println(userID)
	tokenBucket := limiter.store.GetOrCreate(userID)
	bRes := tokenBucket.Consume()

	return LimitResult{
		Allowed:    bRes.Allowed,
		Remaining:  bRes.Remaining,
		RetryAfter: bRes.RetryAfter,
	}, nil
}
