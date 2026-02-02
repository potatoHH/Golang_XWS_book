package repository

import (
	"Book_Exp/webook/interactive/domain"
	"context"
)

type InteractiveRepository interface {
	//阅读数
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	BatchReadCnt(ctx context.Context, biz []string, id []int64) error
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
