package repository

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository/cache"
	cachemocks "Book_Exp/webook/internal/repository/cache/mocks"
	"Book_Exp/webook/internal/repository/dao"
	daomocks "Book_Exp/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	now := time.Now()
	//你要去掉毫秒以外的
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name     string
		ctx      context.Context
		id       int64
		wantErr  error
		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		wantUser domain.User
	}{
		{
			name: "缓存未命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				//缓存未命中
				ud := daomocks.NewMockUserDao(ctrl)
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				ud.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{
					Id:       123,
					Email:    sql.NullString{"123@qq.com", true},
					Password: "123456@www",
					Phone:    sql.NullString{"12345678901", true},
					Ctime:    now.UnixMilli(),
					Utime:    now.UnixMilli(),
				}, nil)
				uc.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456@www",
					Phone:    "12345678901",
					Ctime:    now}).Return(nil)
				return ud, uc
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456@www",
				Phone:    "12345678901",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				//缓存未命中
				ud := daomocks.NewMockUserDao(ctrl)
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456@www",
					Phone:    "12345678901",
					Ctime:    now,
				}, nil)

				return ud, uc
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456@www",
				Phone:    "12345678901",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "数据库查询有误",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				//缓存未命中
				ud := daomocks.NewMockUserDao(ctrl)
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				ud.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{}, errors.New("数据库查询失败"))
				return ud, uc
			},
			ctx:      context.Background(),
			id:       123,
			wantUser: domain.User{},
			wantErr:  errors.New("数据库查询失败"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)

		})
	}

}
