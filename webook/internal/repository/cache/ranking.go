package cache

import (
	"Book_Exp/webook/internal/domain"
	"context"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

type rankingCache struct {
}

func (r *rankingCache) Set(ctx context.Context, arts []domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func (r *rankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func NewRankingCache() RankingCache {
	return &rankingCache{}
}
