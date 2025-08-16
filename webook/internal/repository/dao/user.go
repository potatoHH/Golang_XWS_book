package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}
func (dao *UserDao) Insert(ctx context.Context, u User) error {
	//存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if MysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflicsErrNo uint16 = 1062     //唯一冲突 的错误码
		if MysqlErr.Number == uniqueConflicsErrNo { //是否唯一冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

// 登录
func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error //查询
	return u, err
}
func (dao *UserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error //查询
	return u, err
}

// user在dao层直接对标sql的操作
type User struct {
	Id       int64  `gorm:"primaryKey autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Ctime    int64 // 创建时间
	Utime    int64 //更新时间
}
