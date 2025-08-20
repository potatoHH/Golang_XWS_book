package repository

import (
	"Book_Exp/webook/internal/repository/cache"
	"context"
)

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrInvalidVerifyCode = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputcode string) (bool, error)
}
type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: cache,
	}

}

func (repo *CacheCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}
func (repo *CacheCodeRepository) Verify(ctx context.Context, biz, phone, inputcode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputcode)
}
