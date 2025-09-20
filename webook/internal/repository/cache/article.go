package cache

import (
	"Book_Exp/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, authro int64) error
	Set(ctx context.Context, id int64) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func (r *RedisArticleCache) Set(ctx context.Context, id int64) error {
	data, err := json.Marshal(domain.Article{
		Id: id,
	})
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.Key(id), data, time.Minute*10).Err()
}
func (r *RedisArticleCache) GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error) {
	bs, err := r.client.Get(ctx, r.Key(author)).Bytes()
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	err = json.Unmarshal(bs, &arts)
	return arts, err
}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error {
	for i := range arts {
		//只缓存摘要部分
		arts[i].Content = arts[i].Abstrract()
	}
	bs, err := json.Marshal(arts)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.Key(author),
		bs, time.Minute*10).Err()
}

func (r *RedisArticleCache) Key(uid int64) string {
	return fmt.Sprintf("article:firstpage:%d", uid)
}

func (r *RedisArticleCache) DelFirstPage(ctx context.Context, author int64) error {
	return r.client.Del(ctx, r.Key(author)).Err()
}
