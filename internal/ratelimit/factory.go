package ratelimit

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Factory struct {
	client *redis.Client
}

func NewFactory(client *redis.Client) *Factory {
	return &Factory{client}
}

func (f *Factory) NewLimiter(action string, period time.Duration, limit int64) *Limiter {
	if f == nil {
		return nil
	}
	return NewLimiter(f.client, action, period, limit)
}
