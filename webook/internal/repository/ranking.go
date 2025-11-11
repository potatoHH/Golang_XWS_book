package repository

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository/cache"
	"context"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
}
type CacheRankingRepository struct {
	c cache.RankingCache
}

func (c *CacheRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return c.c.Set(ctx, arts)
}

func NewRankingRepository(c cache.RankingCache) RankingRepository {
	return &CacheRankingRepository{}

}
