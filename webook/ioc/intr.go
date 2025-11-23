package ioc

import (
	intrv1 "Book_Exp/webook/api/proto/gen/intr/v1"
	"Book_Exp/webook/interactive/service"
	"Book_Exp/webook/internal/web/client"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
