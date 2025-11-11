package grpc

import (
	"context"
)

type Server struct {
	UnimplementedUserSeriviceServer
}

func (s Server) GetById(ctx context.Context, request *GetByIdRequest) (*GetByIdResponse, error) {
	return &GetByIdResponse{
		User: &User{
			Id:   345,
			Name: "张三",
		},
	}, nil
}
