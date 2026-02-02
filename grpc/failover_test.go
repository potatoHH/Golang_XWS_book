package grpc

import (
	"Book_Exp/webook/pkg/netx"
	"context"
	_ "embed"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FailoverSuite struct {
	suite.Suite
	client *clientv3.Client
}

func (s *FailoverSuite) SetupSuite() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	require.NoError(s.T(), err)
	s.client = client
}

func (s *FailoverSuite) TestServer() {
	go func() {
		s.startServer(":8091", &AlwaysFailoverServer{
			Name: "AlwaysFailoverServer",
		})
	}()

	s.startServer(":8090", &Server{
		Name: "Server",
	})

}
func (s *FailoverSuite) TestClient() {
	t := s.T()
	cc, err := grpc.NewClient("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NotNil(t, err)
	client := NewUserSeriviceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	resp, err := client.GetById(ctx, &GetByIdRequest{
		Id: 123,
	})
	require.NoError(s.T(), err)
	s.T().Log(resp.User)
}

func (s *FailoverSuite) startServer(addr string, svc UserSeriviceServer) {
	l, err := net.Listen("tcp", addr)
	require.NoError(s.T(), err)
	//endpoint 以服务为维度,一个服务一个Manger
	em, err := endpoints.NewManager(s.client, "service/user")
	require.NoError(s.T(), err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//key是这个实例的key
	//如果 instance id,用 instance id,如果没有 本机iP+端口 进一步区分
	//端口一般是从配置文件里面读取
	addr = netx.GetOutboundIp() + addr
	key := "service/user/" + addr
	//TODO这个ctx是创建租约的超时
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//ttl是租约的租期 长度,秒作为单位
	//过了1/3 还是 3/2的时候开始续约
	var ttl int64 = 30
	leaseResp, err := s.client.Grant(ctx, ttl)
	require.NoError(s.T(), err)
	//TODO 准备发起AddEndpoint
	//AddEndpoint 在这一步之前完成的所有的启动的准备工作,包括缓存预加载之类的事情.
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, clientv3.WithLease(leaseResp.ID))
	require.NoError(s.T(), err)

	//TODO  续约的问题 需要开一个go
	kaCtx, Kacancel := context.WithCancel(context.Background())
	go func() {
		//在这里开一个续约的操作
		ch, err1 := s.client.KeepAlive(kaCtx, leaseResp.ID)
		require.NoError(s.T(), err1)
		for kaResp := range ch {
			//正常就是打印一下DEBUG 日志什么的
			s.T().Log(kaResp.String(), time.Now().String())

		}

	}()

	server := grpc.NewServer()
	RegisterUserSeriviceServer(server, svc)
	err = server.Serve(l)
	s.T().Log(l)
	//正常退出 enpoints
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()
	Kacancel()

	//退出阶段,先在注册中心删除
	err = em.DeleteEndpoint(ctx, key)
	require.NoError(s.T(), err)
	server.GracefulStop() //优雅退出

}

//go:embed failover.json
var svgCfg string

func (s *FailoverSuite) TestRoundRobinClient() {
	db, err := resolver.NewBuilder(s.client)
	require.NoError(s.T(), err)
	cc, err := grpc.Dial("etcd:///service/user/",
		grpc.WithResolvers(db),
		//在这里使用负载均衡
		grpc.WithDefaultServiceConfig(svgCfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := NewUserSeriviceClient(cc)
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.GetById(ctx, &GetByIdRequest{Id: 123})
		cancel()
		require.NoError(s.T(), err)
		s.T().Log(resp)
	}
}

func TestFailoverSuite(t *testing.T) {
	suite.Run(t, new(FailoverSuite))
}
