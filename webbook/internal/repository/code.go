package repository

import (
	"context"
	"jikeshijian_go/webbook/internal/repository/cache"
)

var (
	ErrSetCodeTooManyTimes    = cache.ErrSetCodeTooManyTimes
	ErrCodeVerifyFailed       = cache.ErrCodeVerifyFailed
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type CodeRepositoryWithCache struct {
	codeCache cache.CodeCache
}

func NewCodeRepositoryWithCache(c cache.CodeCache) CodeRepository {
	return &CodeRepositoryWithCache{
		codeCache: c,
	}
}
func (repo *CodeRepositoryWithCache) Store(ctx context.Context, biz string,
	phone string, code string) error {
	err := repo.codeCache.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return nil
}

func (repo *CodeRepositoryWithCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return repo.codeCache.Verify(ctx, biz, phone, inputCode)
}
