package client

import (
	intrv1 "Book_Exp/webook/api/proto/gen/intr/v1"
	"context"
	"math/rand"

	"github.com/ecodeclub/ekit/syncx/atomicx"
	"google.golang.org/grpc"
)

type GreyScaleInteractiveServiceClient struct {
	remote intrv1.InteractiveServiceClient
	local  intrv1.InteractiveServiceClient
	//我们怎么控制流量呢 如果一个请求过来,我该怎么控制他去调用本地,还是调用远程呢,
	//使用随机数+ 阈值的小技巧
	threshold *atomicx.Value[int32]
}

func NewGreyScaleInteractiveServiceClient(remote intrv1.InteractiveServiceClient, local intrv1.InteractiveServiceClient) *GreyScaleInteractiveServiceClient {
	return &GreyScaleInteractiveServiceClient{
		remote:    remote,
		local:     local,
		threshold: atomicx.NewValue[int32](),
	}
}

// 在这里监听startListen 是有缺陷的 GreyScaleInteractiveServiceClient 和viper 是紧耦合
//func (g *GreyScaleInteractiveServiceClient) startLient() {
//	viper.OnConfigChange(func(in fsnotify.Event) {
//
//	})
//}

//// TODO 另一种方法 chan
//func (g *GreyScaleInteractiveServiceClient) Onchange(ch<-chan int32) {
//	go func() {
//		for {
//			NewTh:=range ch{
//				g.UpdateThreshold(NewTh)
//			}
//
//		}
//	}()
//}

func (g *GreyScaleInteractiveServiceClient) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	return g.client().IncrReadCnt(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	return g.client().Like(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	return g.client().CancelLike(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	return g.client().Collect(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	return g.client().Get(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	return g.client().GetByIds(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) UpdateThreshold(newThreshold int32) {
	g.threshold.Store(newThreshold)
}
func (g *GreyScaleInteractiveServiceClient) client() intrv1.InteractiveServiceClient {
	threshold := g.threshold.Load()
	//生成[0-100)随机数
	//举例来说,如果要说 threshold 是100的情况, 可以预见的是所有num 都会进去 返回 g.remote
	num := rand.Int31n(100)
	//如果 threshold 是0 ,那就永远走本地 local
	if num < threshold {
		return g.remote
	}
	return g.local

}
