package ioc

import (
	intrv1 "Book_Exp/webook/api/proto/gen/intr/v1"
	"Book_Exp/webook/interactive/service"
	"Book_Exp/webook/internal/web/client"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 依赖注入式的写法
func InitEtcd() *etcdv3.Client {
	var cfg etcdv3.Config
	err := viper.UnmarshalKey("etcd", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := etcdv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return client
}

// 真正的grpc客户端
func InitGRPCClientV1(client *etcdv3.Client) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr   string
		Secure bool
		Name   string
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	bd, err := resolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{grpc.WithResolvers(bd)}
	if cfg.Secure {
		//上面要去加载你启用的证书之类的
		//启用https
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	//这个地方没有填写正确,也不出错误.
	cc, err := grpc.NewClient("etcd:///service"+cfg.Name, opts...)
	if err != nil {
		panic(err)
	}
	return intrv1.NewInteractiveServiceClient(cc)

}

// 流量控制的客户端吗
func InitGRPCClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string
		Secure    bool
		Threshold int32
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if cfg.Secure {
		//上面要去加载你启用的证书之类的
		//启用https
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.NewClient(cfg.Addr, opts...)
	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewInteractiverServiceApadter(svc)
	res := client.NewGreyScaleInteractiveServiceClient(local, remote)
	//在这里监听配置文件
	viper.OnConfigChange(func(in fsnotify.Event) {
		var cfg Config
		res.UpdateThreshold(cfg.Threshold)
	})
	return res
}
