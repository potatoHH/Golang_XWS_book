package repository

import (
	"Book_Exp/webook/internal/repository/cache"
	"Book_Exp/webook/internal/repository/dao"
	"context"
)

type CacheReadReopsitory interface {
	//点赞
	IncrLike(ctx context.Context, biz string, bizId int64, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId int64, uid int64) error
	//收藏

}
type CachedReadRepository struct {
	dao   dao.InteractiveDAO
	cache cache.RedisInteractiveCache
}

func (c *CachedReadRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	//要考虑缓存方案
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	//这边会有部分数据失败引发的一些不一致问题,但是你其实不需要解决
	//因为阅读数不准确是完全没问题的
	return c.cache.IncrReadCnt(ctx, biz, bizId)
}

func NewInteractiveService() CacheReadReopsitory {
	return &CachedReadRepository{}
}
func (c *CachedReadRepository) DecrLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntPresent(ctx, biz, bizId)
}

func (c *CachedReadRepository) IncrLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.IncrLikeCntPresent(ctx, biz, bizId)
}
