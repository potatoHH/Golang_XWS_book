package repository

import (
	"Book_Exp/webook/internal/repository/cache"
	"context"
)

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrInvalidVerifyCode = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(cache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}

}

func (repo *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}
func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, inputcode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputcode)
}
