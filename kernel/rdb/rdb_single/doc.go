package rdb_single

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var defRedis *RedisSingle

func InitSingle(ctx context.Context, opts ...any) (err error) {
	defRedis = NewRedisSingle()
	if err = defRedis.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *redis.Client {
	return defRedis.Get().(*redis.Client)
}

func GetCtx() context.Context {
	return defRedis.GetCtx()
}
