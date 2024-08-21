package rdb_single

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisSingle struct {
	options *RedisSingleOption
	client  *redis.Client

	ctx    context.Context
	cancel context.CancelFunc
}

func NewRedisSingle() *RedisSingle {
	rs := &RedisSingle{
		options: NewRedisSingleOption(),
	}

	return rs
}

func (r *RedisSingle) Init(ctx context.Context, opts ...any) (err error) {
	r.ctx, r.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(r.options)
		}
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:     r.options.addr,
		Password: r.options.pwd,
		DB:       0,
		PoolSize: 300,
	})
	if err = r.client.Ping(r.ctx).Err(); err != nil {
		return
	}

	return nil

}

func (r *RedisSingle) Get() any {
	return r.client
}

func (r *RedisSingle) GetCtx() context.Context {
	return r.ctx
}
