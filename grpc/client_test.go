package grpc

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClient(t *testing.T) {
	//cc就是一个连接池,cc里面放了好多个连接
	cc, err := grpc.NewClient(":8090",
		grpc.WithTransportCredentials(
			insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(firstClient, secondClient))
	require.NoError(t, err)
	client := NewUserSeriviceClient(cc)
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*30)
	defer cancle()
	resp, err := client.GetById(ctx, &GetByIdRequest{
		Id: 345,
	})
	assert.NoError(t, err)
	t.Log(resp.User)
}

var firstClient = grpc.UnaryClientInterceptor(func(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Println("客户端第一个")
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Println("客户端第一个后")
	return err
})
var secondClient = grpc.UnaryClientInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Println("客户端第二个前")
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Println("客户端第二个后")
	return err

})
