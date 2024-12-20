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

type CodeRepository struct {
	codeCache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		codeCache: c,
	}
}
func (repo *CodeRepository) Store(ctx context.Context, biz string,
	phone string, code string) error {
	err := repo.codeCache.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return nil
}

func (repo *CodeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return repo.codeCache.Verify(ctx, biz, phone, inputCode)
}
