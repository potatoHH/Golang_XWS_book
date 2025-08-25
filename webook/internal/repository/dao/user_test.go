package dao

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGormUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		ctx     context.Context
		mock    func(t *testing.T) *sql.DB //因为这里不是gomock，所以这里需要自己模拟数据库
		user    User
		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				//这边预期的是正则表达式
				mock.ExpectExec("INSERT INTO `users`.*").
					WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				//这边预期的是正则表达式
				mock.ExpectExec("INSERT INTO `users`.*").
					WillReturnError(&mysql.MySQLError{
						Number: 1062,
					})
				require.NoError(t, err)
				return mockDB
			},
			user:    User{},
			wantErr: ErrUserDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				//res := sqlmock.NewResult(3, 1)
				//这边预期的是正则表达式
				mock.ExpectExec("INSERT INTO `users`.*").
					WillReturnError(errors.New("数据库错误"))
				require.NoError(t, err)
				return mockDB
			},
			user:    User{},
			wantErr: errors.New("数据库错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(gormmysql.New(gormmysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true, // 跳过版本检测
			}), &gorm.Config{
				DisableAutomaticPing:   true, // 关闭自动ping
				SkipDefaultTransaction: true, // 关闭默认事务
			})
			d := NewUserDao(db)
			u := tc.user
			err = d.Insert(tc.ctx, u)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
