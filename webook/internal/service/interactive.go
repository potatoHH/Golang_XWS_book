package service

import (
	"Book_Exp/webook/internal/repository"
	"context"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	//点赞
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	CancleLike(ctx context.Context, biz string, bizId int64, uid int64) error
}
type InteractiveServiceV1 struct {
	repo repository.InteractiveService
}

func (a *InteractiveServiceV1) CancleLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return a.repo.DecrLike(ctx, biz, bizId, uid)
}

func (a *InteractiveServiceV1) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	return a.repo.IncrLike(ctx, biz, bizId, uid)
}

func NewInteractiveService() InteractiveService {
	return &InteractiveServiceV1{}
}
func (i *InteractiveServiceV1) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}
