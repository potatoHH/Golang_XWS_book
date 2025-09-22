package repository

import (
	"Book_Exp/webook/internal/repository/cache"
	"Book_Exp/webook/internal/repository/dao"
	"context"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
}

type InteractiveServiceV1 struct {
	dao   dao.InteractiveDAO
	cache cache.RedisInteractiveCache
}

func (i InteractiveServiceV1) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	//要考虑缓存方案
	err := i.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	//这边会有部分数据失败引发的一些不一致问题,但是你其实不需要解决
	//因为阅读数不准确是完全没问题的
	return i.cache.IncrReadCnt(ctx, biz, bizId)
}

func NewInteractiveService() InteractiveService {
	return &InteractiveServiceV1{}
}
