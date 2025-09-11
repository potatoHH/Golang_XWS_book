package web

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/service"
	svcmocks "Book_Exp/webook/internal/service/mocks"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// 密码加密测试
func TestPasswordEncrypt(t *testing.T) {
	password := "123456"
	encrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypt, []byte(password))
	assert.NoError(t, err) // 断言
}

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserServiceV1
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserServiceV1 {
				usersvc := svcmocks.NewMockUserServiceV1(ctrl)
				//codesvc := svcmocks.NewMockCodeServiceV1(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "2138824181@qq.com",
					Password: "123456@ww",
				}).Return(nil)
				return usersvc
			},
			reqBody: `{
					"email":"2138824181@qq.com",
					"password":"123456@ww",
					"confirmPassword":"123456@ww",
				}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对,bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserServiceV1 {
				usersvc := svcmocks.NewMockUserServiceV1(ctrl)
				//codesvc := svcmocks.NewMockCodeServiceV1(ctrl)
				return usersvc
			},
			reqBody: `{
					"email":"2138824181@qq.com",
					"password":"123456@ww"
				}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserServiceV1 {
				usersvc := svcmocks.NewMockUserServiceV1(ctrl)
				//codesvc := svcmocks.NewMockCodeServiceV1(ctrl)
				return usersvc
			},
			reqBody: `{
					"email":"2138824181",
					"password":"123456@ww"
					"confirmPassword":"123456@ww",
				}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱格式不对",
		},
		{
			name: "密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserServiceV1 {
				usersvc := svcmocks.NewMockUserServiceV1(ctrl)
				//codesvc := svcmocks.NewMockCodeServiceV1(ctrl)
				return usersvc
			},
			reqBody: `{
					"email":"2138824181",
					"password":"123456@ww"
					"confirmPassword":"123456@w",
				}`,
			wantCode: http.StatusOK,
			wantBody: "密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserServiceV1 {
				usersvc := svcmocks.NewMockUserServiceV1(ctrl)
				//codesvc := svcmocks.NewMockCodeServiceV1(ctrl)
				return usersvc
			},
			reqBody: `{
					"email":"2138824181",
					"password":"123456"
					"confirmPassword":"123456",
				}`,
			wantCode: http.StatusOK,
			wantBody: "密码格式错误",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserServiceV1 {
				usersvc := svcmocks.NewMockUserServiceV1(ctrl)
				//codesvc := svcmocks.NewMockCodeServiceV1(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "2138824181@qq.com",
					Password: "123456@ww",
				}).Return(service.ErrUserDuplicate)
				return usersvc
			},
			reqBody: `{
					"email":"2138824181@qq.com",
					"password":"123456@ww",
					"confirmPassword":"123456@ww",
				}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserServiceV1 {
				usersvc := svcmocks.NewMockUserServiceV1(ctrl)
				//codesvc := svcmocks.NewMockCodeServiceV1(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "2138824181@qq.com",
					Password: "123456@ww",
				}).Return(errors.New("error 系统"))
				return usersvc
			},
			reqBody: `{
					"email":"2138824181@qq.com",
					"password":"123456@ww",
					"confirmPassword":"123456@ww",
				}`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}
	ring())
		})
	}
}for _, tc := range testCases {
t.Run(tc.name, func(t *testing.T) {
ctrl := gomock.NewController(t)
defer ctrl.Finish()
//TODO注册路由
server := gin.Default()
//要用mock模拟,这里还没有处理
hdl := NewUserHandler(tc.mock(ctrl), nil)
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

func TestMock(t *testing.T) {
	//先创建一个控制mock 的控制器
	ctrl := gomock.NewController(t)
	//每一个测试结束都要调用finish
	//然后mock 就会验证你的测试流程是否符合预期
	defer ctrl.Finish()
	usersvc := svcmocks.NewMockUserServiceV1(ctrl)
	//开始设计一个模拟调用
	//预期的哥是Signup 的调用
	//模拟的条件是gomock.Any,gomock.Any. 随便传 .Any()
	//然后返回调用
	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
		Return(errors.New("mock err"))
	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})
	t.Log(err)
}
