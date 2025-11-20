package main

import (
	intrv1 "Book_Exp/webook/api/proto/gen/intr/v1"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPClient(t *testing.T) {
	// 连接服务端
	cc, err := grpc.NewClient("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := intrv1.NewInteractiveServiceClient(cc)
	// 调用 RPC
	resp, err := client.Get(context.Background(), &intrv1.GetRequest{
		Biz:   "test",
		BizId: 2,
		Uid:   456,
	})
	require.NoError(t, err)
	t.Log(resp.Intr)
}
