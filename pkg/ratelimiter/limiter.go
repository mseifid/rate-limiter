package ratelimiter

import (
	"context"
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
	//userID := ctx.Value("userID")
	tokenBucket := limiter.store.GetOrCreate("1")
	bRes := tokenBucket.Consume()

	return LimitResult{
		Allowed:    bRes.Allowed,
		Remaining:  bRes.Remaining,
		RetryAfter: bRes.RetryAfter,
	}, nil
}
