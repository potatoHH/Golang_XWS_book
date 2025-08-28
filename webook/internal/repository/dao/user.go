package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicate = errors.New("邮箱冲突 or 手机号冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	UpdateNonZeroFields(ctx context.Context, u User) error
	FindByWecaht(ctx context.Context, openId string) (User, error)
}

type GormUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GormUserDao{
		db: db,
	}
}
func (dao *GormUserDao) Insert(ctx context.Context, u User) error {
	//存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if MysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflicsErrNo uint16 = 1062     //唯一冲突 的错误码
		if MysqlErr.Number == uniqueConflicsErrNo { //是否唯一冲突
			//邮箱  or 手机号码 冲突
			return ErrUserDuplicate
		}
	}
	return err
}

// 登录
func (dao *GormUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error //查询
	return u, err
}
func (dao *GormUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error //查询
	return u, err
}
func (dao *GormUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error //查询
	return u, err
}

func (ud *GormUserDao) UpdateNonZeroFields(ctx context.Context, u User) error {
	// 这种写法是很不清晰的，因为它依赖了 gorm 的两个默认语义
	// 会使用 ID 来作为 WHERE 条件
	// 会使用非零值来更新
	// 另外一种做法是显式指定只更新必要的字段，
	// 那么这意味着 DAO 和 failover 中非敏感字段语义耦合了
	return ud.db.Updates(&u).Error
}

// user在dao层直接对标sql的操作
type User struct {
	Id int64 `gorm:"primaryKey autoIncrement"`
	//这两个字段会造成索引冲突  用到sql.NullStrin  允许有多个空值  但不允许有多个""
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"` //phone *string  这种写法你要解引用 ,你要判空
	Password string
	Ctime    int64 // 创建时间
	Utime    int64 //更新时间
	//微信的字段
	WechatUnionId sql.NullString `gorm:"unique"`
	WechatOpenId  sql.NullString `gorm:"unique"`
}
