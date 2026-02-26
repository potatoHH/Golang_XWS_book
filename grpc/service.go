package grpc

import (
	"context"
	"log"
	"time"
)

type Server struct {
	UnimplementedUserServiceServer
	Name string
}

func (s *Server) GetById(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error) {
	//链路超时
	ddl, ok := ctx.Deadline()
	if ok {
		//打印剩余超时时间
		log.Println(ddl.Sub(time.Now()).String())
	}

	return &GetByIdResponse{
		User: &User{
			Id:   req.Id,
			Name: s.Name,
		},
	}, nil
}
