package repository

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository/cache"
	"Book_Exp/webook/internal/repository/dao"
	"Book_Exp/webook/pkg/logger"
	"context"
)

type CacheReadReopsitory interface {
	//阅读数
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	//点赞
	IncrLike(ctx context.Context, biz string, bizId int64, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId int64, uid int64) error
	//收藏
	AddCollectItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error)
	//点过赞/收藏
	Liked(ctx context.Context, biz string, bizId, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, bizId, uid int64) (bool, error)
}
type CachedReadRepository struct {
	dao   dao.InteractiveDAO
	cache cache.RedisInteractiveCache
	l     logger.LoggerV1
}

func (c *CachedReadRepository) Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, bizId)

	if err == nil {
		// 缓存只缓存了具体的数字，但是没有缓存自身有没有点赞的信息
		// 因为一个人反复刷，重复刷一篇文章是小概率的事情
		// 也就是说，你缓存了某个用户是否点赞的数据，命中率会很低
		return intr, nil
	}
	daoInt, err := c.dao.Get(ctx, biz, bizId)
	if err == dao.ErrRecordNotFound || err == nil {
		res := c.toDomain(daoInt)
		go func() {
			if er := c.cache.Set(ctx, biz, bizId, res); er != nil {
				c.l.Error("回写缓存失败",
					logger.Int64("bizId", bizId),
					logger.String("biz", biz),
					logger.Error(er))
			}
		}()
		return res, nil
	}
	return domain.Interactive{}, err
}

func (c *CachedReadRepository) Liked(ctx context.Context, biz string, bizId, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, bizId, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedReadRepository) Collected(ctx context.Context, biz string, bizId, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, bizId, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedReadRepository) AddCollectItem(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error {
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		BizId: bizId,
		Cid:   cid,
		Uid:   uid,
	})
	//这个地方,你要不要考虑缓存 以及收藏夹的内内容,用户会频繁访问它的收藏夹,那么你就应该缓存,不然你就不需要
	if err != nil {
		return err
	}
	return c.cache.IncrCollectCntIfPresent(ctx, biz, bizId)

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

func NewInteractiveService(dao dao.InteractiveDAO, cache cache.RedisInteractiveCache, l logger.LoggerV1) CacheReadReopsitory {
	return &CachedReadRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
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

// TODO 正常来说,参数必然不用指针:方法不要修改参数,通过返回值来修改参数
// TODO 返回值就看情况,如果是指针实现了接口,那么就返回指针,如果返回值很大,你不想传递引发复制问题,那么还是返回指针
func (c *CachedReadRepository) toDomain(intr dao.Interactive) domain.Interactive {
	return domain.Interactive{
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		ReadCnt:    intr.ReadCnt,
	}
}

// TODO 不好的实践和凑活的实践
func (c *CachedReadRepository) UpdateCntV1(intr *dao.Interactive) {
	intr.LikeCnt = 30
}
func (c *CachedReadRepository) UpdateCntV2(intr dao.Interactive) domain.Interactive {
	intr.LikeCnt = 30
	return c.toDomain(intr)
}
