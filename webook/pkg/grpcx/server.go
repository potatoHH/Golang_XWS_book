package grpcx

import (
	"Book_Exp/webook/pkg/logger"
	"Book_Exp/webook/pkg/netx"
	"context"
	"net"
	"strconv"
	"time"

	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Port     int
	EtcdAddr []string
	L        logger.LoggerV1
	client   *etcdv3.Client
	Name     string
	em       endpoints.Manager
	kaCancel func()
	key      string
}

func NewServer(client etcdv3.Client) {

}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		return err
	}

	//TODO  在这直接启动,要在这里创建服务注册的过程

	//这边会阻塞,类似与gin.Run
	return s.Server.Serve(l)

}
func (s *Server) Register() error {
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints: s.EtcdAddr,
	})
	if err != nil {
		return err
	}
	err = s.Register()
	if err != nil {
		s.L.Error("连接服务失败", logger.Error(err))
	}
	//endpoint 以服务为维度,一个服务一个Manger
	em, err := endpoints.NewManager(client, "service/"+s.Name)
	if err != nil {
		return err
	}
	s.client = client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	addr := netx.GetOutboundIp() + ":" + strconv.Itoa(s.Port) //
	key := "service/" + s.Name + "/" + addr
	s.key = key
	//TODO这个ctx是创建租约的超时
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//ttl是租约的租期 长度,秒作为单位
	//过了1/3 还是 3/2的时候开始续约
	var ttl int64 = 30
	leaseResp, err := client.Grant(ctx, ttl)
	if err != nil {
		return err
	}
	//TODO 准备发起AddEndpoint
	//AddEndpoint 在这一步之前完成的所有的启动的准备工作,包括缓存预加载之类的事情.
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, etcdv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	//TODO  续约的问题 需要卡一个go
	kaCtx, Kacancel := context.WithCancel(context.Background())
	s.kaCancel = Kacancel
	//在这里开一个续约的操作
	ch, err1 := client.KeepAlive(kaCtx, leaseResp.ID)
	if err1 != nil {
		logger.L().Error("租约过期", logger.Error(err1))
	}
	go func() {
		for kaResp := range ch {
			//正常就是打印一下DEBUG 日志什么的
			s.L.Debug(kaResp.String())
		}
	}()
	return nil

}

// Close 可以叫做shutdown
func (s *Server) Close() error {
	if s.kaCancel != nil {
		s.kaCancel()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if s.em != nil {
		err := s.em.DeleteEndpoint(ctx, s.key)
		if err != nil {
			return err
		}
	}
	if s.client != nil {
		err := s.client.Close()
		if err != nil {
			return err
		}
	}
	//优雅退出
	s.Server.GracefulStop()
	return nil
}
