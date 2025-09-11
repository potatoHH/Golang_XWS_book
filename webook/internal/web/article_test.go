package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestArticlehandler_Publish(t *testing.T){
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			//TODO注册路由
			server := gin.Default()
			//要用mock模拟,这里还没有处理
			hdl := NewUserHandler(tc.mock(ctrl), nil,nil)
			hdl.RegisterRoutes(server)
			//TODO 构造请求
			req, err := http.NewRequest(http.MethodPost, "/users/signup",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err) //要求一定没有err
			//数据是json格式
			req.Header.Set("Content-Type", "application/json")
			t.Log(req)
			//这里就正常使用req
			resp := httptest.NewRecorder()
			//请求gin的路口  响应返回到resp
			server.ServeHTTP(resp, req)
			//断言结果
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.St
}
