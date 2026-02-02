package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AlwaysFailoverServer struct {
	UnimplementedUserSeriviceServer
	Name string
}

var _ UserSeriviceServer = &Server{}

func (s *AlwaysFailoverServer) GetById(ctx context.Context, request *GetByIdRequest) (*GetByIdResponse, error) {
	fmt.Println("进入了fail服务端")
	return &GetByIdResponse{
		User: &User{
			Id:   1233,
			Name: "来自永远失败的节点" + s.Name,
		},
	}, status.Errorf(codes.Unavailable, "模拟服务端异常")
}
