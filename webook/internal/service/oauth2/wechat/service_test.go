//go:build e2e

package wechat

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_service_e2e_VerityCode(t *testing.T) {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_ID")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_SECRET")
	}
	svc := NewService(appId, appKey,nil)
	res, err := svc.VerityCode(context.Background(), "code", "state")
	require.NoError(t, err)
	t.Log(res)
}
