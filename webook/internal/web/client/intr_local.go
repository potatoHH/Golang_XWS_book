package client

import (
	intrv1 "Book_Exp/webook/api/proto/gen/intr/v1"
	"Book_Exp/webook/interactive/domain"
	"Book_Exp/webook/interactive/service"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//将一个本地实现来伪装成grpc客户端

type InteractiverServiceApadter struct {
	svc service.InteractiveService
}

func NewInteractiverServiceApadter(svc service.InteractiveService) *InteractiverServiceApadter {
	return &InteractiverServiceApadter{svc: svc}
}

func (i *InteractiverServiceApadter) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, in.Biz, in.BizId)
	return &intrv1.IncrReadCntResponse{}, err
}

func (i *InteractiverServiceApadter) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	err := i.svc.Like(ctx, in.Biz, in.BizId, in.Uid)
	return &intrv1.LikeResponse{}, err
}

func (i *InteractiverServiceApadter) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	err := i.svc.CancleLike(ctx, in.Biz, in.BizId, in.Uid)
	return &intrv1.CancelLikeResponse{}, err
}

func (i *InteractiverServiceApadter) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	if in.Uid <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Uid 非法")
	}
	err := i.svc.Collect(ctx, in.Biz, in.BizId, in.Uid, in.Cid)
	return &intrv1.CollectResponse{}, err
}

func (i *InteractiverServiceApadter) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	res, err := i.svc.Get(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Intr: i.toDTO(res),
	}, nil
}
func (i *InteractiverServiceApadter) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	_, err := i.svc.GetByIds(ctx, in.Biz, in.Ids)
	return &intrv1.GetByIdsResponse{}, err
}

func (i *InteractiverServiceApadter) toDTO(intr domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collected,
	}
}
