package ratelimit

import "context"

type Limiter interface {
	// Limited 有没有出发限流 . key是限流对象
	// bool 是否限流
	// 限流器本身是否存在错误
	Limited(ctx context.Context, key string) (bool, error)
}
