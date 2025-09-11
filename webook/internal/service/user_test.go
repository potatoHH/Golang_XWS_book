package service

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

//	func TestUserService_Login(t *testing.T) {
//		now := time.Now() //获取当前时间
//		testCases := []struct {
//			name string
//			mock func(ctrl *gomock.Controller) repository.UserRepository
//			//输入
//			ctx      context.Context
//			email    string
//			password string
//			wantErrr error
//			wantUser domain.User
//		}{
//			{
//				name: "登录成功",
//				mock: func(ctrl *gomock.Controller) repository.UserRepository {
//					repo := repomocks.NewMockUserRepository(ctrl)
//					repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
//						Return(domain.User{
//							Email:    "123@qq.com",
//							Password: "$2a$10$zX6ft5/wGSYBQiAnbcJM.O.tzONexxYhw4cVVrhqnDqbso5h4Lxe.",
//							Phone:    "15231221452",
//							Ctime:    now,
//						}, nil)
//					return repo
//				},
//				email:    "123@qq.com",
//				password: "123456@www",
//
//				wantUser: domain.User{
//					Email:    "123@qq.com",
//					Password: "$2a$10$zX6ft5/wGSYBQiAnbcJM.O.tzONexxYhw4cVVrhqnDqbso5h4Lxe.",
//					Phone:    "15231221452",
//					Ctime:    now,
//				},
//				wantErrr: nil,
//			},
//			{
//				name: "用户不存在",
//				mock: func(ctrl *gomock.Controller) repository.UserRepository {
//					repo := repomocks.NewMockUserRepository(ctrl)
//					repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
//						Return(domain.User{}, repository.ErrUserNotFound)
//					return repo
//				},
//				email:    "123@qq.com",
//				password: "123456@www",
//
//				wantUser: domain.User{},
//				wantErrr: ErrInvalidUserOrPassword,
//			},
//			{
//				name: "系统错误",
//				mock: func(ctrl *gomock.Controller) repository.UserRepository {
//					repo := repomocks.NewMockUserRepository(ctrl)
//					repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
//						Return(domain.User{}, errors.New("系统错误"))
//					return repo
//				},
//				email:    "123@qq.com",
//				password: "123456@www",
//
//				wantUser: domain.User{},
//				wantErrr: errors.New("系统错误"),
//			},
//			{
//				name: "密码错误",
//				mock: func(ctrl *gomock.Controller) repository.UserRepository {
//					repo := repomocks.NewMockUserRepository(ctrl)
//					repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
//						Return(domain.User{
//							Email:    "123@qq.com",
//							Password: "$2a$10$zX6ft5/wGSYBQiAnbcJM.O.tzONexxYhw4cVVrhqnDqbso5h4Lxe.",
//							Phone:    "15231221452",
//							Ctime:    now,
//						}, nil)
//					return repo
//				},
//				email:    "123@qq.com",
//				password: "123456@www1",
//
//				wantUser: domain.User{},
//				wantErrr: ErrInvalidUserOrPassword,
//			},
//		}
//		for _, tc := range testCases {
//			t.Run(tc.name, func(t *testing.T) {
//				//具体的测试代码
//				ctrl := gomock.NewController(t)
//				defer ctrl.Finish()
//				svc := NewUserService(tc.mock(ctrl))
//				u, err := svc.Login(tc.ctx, tc.email, tc.password)
//				assert.Equal(t, tc.wantErrr, err)
//				assert.Equal(t, tc.wantUser, u)
//
//			})
//
//		}
//	}
func TestEncrypt(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("123456@www"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
