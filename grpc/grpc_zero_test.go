package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type GoZeroTest struct {
	suite.Suite
	client *etcdv3.Client
}

func (s *GoZeroTest) TestGoZeroServer() {
	c := zrpc.RpcServerConf{
		//这个是服务启动的地址
		ListenOn: "8090",
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:12379"},
			Key:   "user",
		},
	}
	server := zrpc.MustNewServer(c, func(server *grpc.Server) {
		//把你的业务注册到你的server里
		RegisterUserServiceServer(server, &Server{})
	})
	server.Start()
}
func (s *GoZeroTest) TestGoZeroClient() {
	zClient := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:12379"},
			Key:   "user",
		},
	})
	client := NewUserServiceClient(zClient.Conn())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.GetById(ctx, &GetByIdRequest{Id: 1})
	require.NoError(s.T(), err)
	s.T().Log(resp)
}

func TestGoZeroTestSuite(t *testing.T) {
	suite.Run(t, new(GoZeroTest))

}
