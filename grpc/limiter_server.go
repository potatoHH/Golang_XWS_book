package grpc

import (
	"Book_Exp/webook/pkg/logger"
	"Book_Exp/webook/pkg/ratelimit"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LimiterUserServer struct {
	limiter ratelimit.Limiter
	UserServiceServer
	l logger.LoggerV1
}

func (s *LimiterUserServer) GetById(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error) {
	limited, err := s.limiter.Limit(ctx, fmt.Sprintf("limiter:user:get_by_Id:%d", req.Id))
	if err != nil {
		s.l.Error("判定限流出现问题", logger.Error(err))
		return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
	}
	resp, err := s.UserServiceServer.GetById(ctx, req)
	if limited {
		return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
	}
	return resp, nil
}
func (s *LimiterUserServer) UpdateById(ctx context.Context, req *UpdateRequest) (*UpdateByIdResponse, error) {
	limited, err := s.limiter.Limit(ctx, fmt.Sprintf("limiter:user:update_by_Id:%d", req.Id))
	if err != nil {
		s.l.Error("判定限流出现问题", logger.Error(err))
		return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
	}
	if limited {
		return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
	}

	resp, err := s.UserServiceServer.UpdateById(ctx, req)
	if err != nil {
		s.l.Error("判定限流出现问题", logger.Error(err))
		return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
	}
	return resp, nil

}

func (s *LimiterUserServer) BuilderServerInterceptorBiz() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		if getId, ok := req.(*GetByIdRequest); ok {
			key := fmt.Sprintf("limiter:user:get_by_id_:%d", getId.Id)
			ok, err = s.limiter.Limit(ctx, key)
			if err != nil {
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}
			if ok {
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}
		}

		return handler(ctx, req)
	}

}
