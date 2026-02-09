package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClient(t *testing.T) {
	//cc就是一个连接池,cc里面放了好多个连接
	cc, err := grpc.NewClient(":8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
