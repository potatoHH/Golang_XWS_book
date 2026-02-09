package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
)

func TestServer(t *testing.T) {
	//TODO  创建一个grpc的server
	server := grpc.NewServer()
	defer func() {
		//优雅退出
		server.GracefulStop()
	}()
	//TODO 创建一个user的server
	userServer := &Server{}
	RegisterUserSeriviceServer(server, userServer)
	//TODO 创建一个监听器,监听tcp端口8090
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	//TODO 启动server
	server.Serve(listener)
	t.Log(err)

}

// 服务治理 -服务级别限流
func TestGrpcServer(t *testing.T) {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(first, second))
	defer func() {
		//优雅退出
		server.GracefulStop()
	}()
	userServer := &Server{}
	RegisterUserSeriviceServer(server, userServer)
	//TODO 创建一个监听器,监听tcp端口8090
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	//TODO 启动server
	server.Serve(listener)
	t.Log(err)

}

var first grpc.UnaryServerInterceptor = func(ctx context.Context,
	req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("这是第一个前")
	resp, err = handler(ctx, req)
	log.Println("这是第一个后")
	return resp, err
}
var second grpc.UnaryServerInterceptor = func(ctx context.Context,
	req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("这是第二个前")
	resp, err = handler(ctx, req)
	log.Println("这是第二个后")
	return resp, err
}
