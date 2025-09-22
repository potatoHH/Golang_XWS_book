package service

import (
	"Book_Exp/webook/internal/repository"
	"context"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
}
type InteractiveServiceV1 struct {
	repo repository.InteractiveService
}

func NewInteractiveService() InteractiveService {
	return &InteractiveServiceV1{}
}
func (i *InteractiveServiceV1) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}
