package service

import (
	"Book_Exp/webook/interactive/domain"
	"Book_Exp/webook/interactive/repository"
	"Book_Exp/webook/pkg/logger"
	"context"

	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -source=./interactive.go -package=svcmocks -destination=mocks/interactive.mock.go InteractiveService
type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	//点赞
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	CancleLike(ctx context.Context, biz string, bizId int64, uid int64) error
	//收藏 cid 为收藏夹的id
	Collect(ctx context.Context, biz string, bizId, cid, uid int64) error
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	GetByIds(ctx context.Context, biz string, bizids []int64) (map[int64]domain.Interactive, error)
}
type InteractiveServiceV1 struct {
	repo repository.InteractiveRepository
	l    logger.LoggerV1
}

func (a *InteractiveServiceV1) GetByIds(ctx context.Context, biz string, bizids []int64) (map[int64]domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (a *InteractiveServiceV1) Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error) {
	//TODO 按照repostory 的语义,你这里拿到的应该是全部的数据
	var (
		eg        errgroup.Group
		intr      domain.Interactive
		collected bool
		liked     bool
	)
	eg.Go(func() error {
		var err error
		intr, err = a.repo.Get(ctx, biz, bizId, uid)
		return err
	})
	eg.Go(func() error {
		var err error
		liked, err = a.repo.Liked(ctx, biz, bizId, uid)
		return err
	})
	eg.Go(func() error {
		var err error
		collected, err = a.repo.Collected(ctx, biz, bizId, uid)
		return err
	})
	err := eg.Wait()
	if err != nil {
		return domain.Interactive{}, err
	}
	intr.Liked = liked
	intr.Collected = collected
	return intr, nil
}

func (a *InteractiveServiceV1) Collect(ctx context.Context, biz string, bizId, cid, uid int64) error {
	return a.repo.AddCollectItem(ctx, biz, bizId, cid, uid)
}

func (a *InteractiveServiceV1) CancleLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return a.repo.DecrLike(ctx, biz, bizId, uid)
}

func (a *InteractiveServiceV1) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	return a.repo.IncrLike(ctx, biz, bizId, uid)
}

func NewInteractiveService(repo repository.InteractiveRepository, l logger.LoggerV1) InteractiveService {
	return &InteractiveServiceV1{
		repo: repo,
		l:    l,
	}
}
func (a *InteractiveServiceV1) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return a.repo.IncrReadCnt(ctx, biz, bizId)
}
