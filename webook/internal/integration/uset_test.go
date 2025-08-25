package integration

import (
	"Book_Exp/webook/internal/web"
	"Book_Exp/webook/ioc"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_SendLoginSmsCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		//TODO 你要考虑准备数据和验证数据  slq 和redis 的数据对不对
		reqBody  string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				//Redis没有任何数据
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				//你要清理数据
				val, err := (rdb.GetDel(ctx, "phone_code:login:15321551609").Result())
				cancel()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6) //验证码为6位

			},
			reqBody: `{
				"phone":"15321551609"
				}`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				//Redis没有任何数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				//你要清理数据
				_, err := rdb.Set(ctx, "phone_code:login:15321551609", "123456", time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				//你要清理数据
				val, err := rdb.GetDel(ctx, "phone_code:login:15321551609").Result()
				cancel()
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)

			},
			reqBody: `{
				"phone":"15321551609"
				}`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "验证码发送次数太多,请稍后再试",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				//Redis没有任何数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				//你要清理数据
				_, err := rdb.Set(ctx, "phone_code:login:15321551609", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				//你要清理数据
				val, err := rdb.GetDel(ctx, "phone_code:login:15321551609").Result()
				cancel()
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)

			},
			reqBody: `{
				"phone":"15321551609"
				}`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "手机号为空",
			before: func(t *testing.T) {
				//Redis没有任何数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				//你要清理数据
				_, err := rdb.Set(ctx, "phone_code:login:15321551609", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				//你要清理数据
				val, err := rdb.GetDel(ctx, "phone_code:login:15321551609").Result()
				cancel()
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)

			},
			reqBody: `{
				"phone":""
				}`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 4,
				Msg:  "输入有误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			//TODO 构造请求
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err) //要求一定没有err
			//数据是json格式
			req.Header.Set("Content-Type", "application/json")
			//这里就正常使用req
			resp := httptest.NewRecorder()
			//请求gin的路口  响应返回到resp
			server.ServeHTTP(resp, req)
			//断言结果
			assert.Equal(t, tc.wantCode, resp.Code)
			var webRes web.Result
			json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webRes)
			tc.after(t)

		})
	}

}
