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
	GetPublishedById(ctx context.Context, id int64) ([]domain.Article, error)
	// SetPub 正常来说看,创作者和读者的Reids 集成要分开,因为读者是一个核心中的核心
	SetPub(ctx context.Context, article domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	//点赞
	IncrLikeCntPresent(ctx context.Context, biz string, id int64) error
	DecrLikeCntPresent(ctx context.Context, biz string, id int64) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func (r *RedisArticleCache) IncrLilkeCntPresent(ctx context.Context, biz string, id int64) {}

func (r *RedisArticleCache) Set(ctx context.Context, id int64) error {
	data, err := json.Marshal(domain.Article{
		Id: id,
	})
	if err != nil {
		return err
	}
	//TODO 过期时间要段.     你的预测效果越不好,就越要短
	return r.client.Set(ctx, r.Key(id), data, time.Minute*10).Err()
}
func (r *RedisArticleCache) GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error) {
	bs, err := r.client.Get(ctx, r.firstPageKey(author)).Bytes()
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

	return r.client.Set(ctx, r.firstPageKey(author),
		bs, time.Minute*10).Err()
}

func (r *RedisArticleCache) firstPageKey(uid int64) string {
	return fmt.Sprintf("article:firstpage:%d", uid)
}
func (r *RedisArticleCache) Key(id int64) string {
	return fmt.Sprintf("article:%d", id)
}

func (r *RedisArticleCache) DelFirstPage(ctx context.Context, author int64) error {
	return r.client.Del(ctx, r.firstPageKey(author)).Err()
}
