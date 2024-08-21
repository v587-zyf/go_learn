package iface

import "context"

type IRedis interface {
	Init(ctx context.Context, opts ...any) (err error)
	Get() any
	GetCtx() context.Context
}
