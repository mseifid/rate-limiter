package ratelimiter

import (
	"context"
)

type Limiter struct {
	store Store
}

type Store interface {
	GetOrCreate(userID string) *Bucket
	GetGlobal() *Bucket
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
	globalBucket := limiter.store.GetGlobal()
	guRes := globalBucket.Consume()
	if !guRes.Allowed {
		return LimitResult{
			Allowed:    guRes.Allowed,
			Remaining:  guRes.Remaining,
			RetryAfter: guRes.RetryAfter,
		}, nil
	}

	userID := ctx.Value("userID").(string)
	userBucket := limiter.store.GetOrCreate(userID)
	ubRes := userBucket.Consume()

	return LimitResult{
		Allowed:    ubRes.Allowed,
		Remaining:  ubRes.Remaining,
		RetryAfter: ubRes.RetryAfter,
	}, nil
}
