package grpc

import (
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
