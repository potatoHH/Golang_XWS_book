package repository

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository/cache"
	"Book_Exp/webook/internal/repository/dao"
	"context"
	"database/sql"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	Update(ctx context.Context, u domain.User) error
}

type CachedUserRepository struct { // 用户仓库
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}
func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}
func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	//缓存里面有数据,缓存里面没有数据,缓存出错
	if err == nil {
		//必然有数据
		return u, nil
	}
	//缓存中没有数据 ,去数据库里面找
	//if err != cache.ErrKeyNotExist {
	//	//去数据库里加载
	//}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = r.entityToDomain(ue)
	_ = r.cache.Set(ctx, u)
	//go func() {
	//	//打日志,做监控
	//	if err != nil {
	//
	//	}
	//}()
	return u, nil

	//这里怎么版?err =io.EOF
	//要不要去数据库加载?
	//看起来我不应该加载
	//看起来我好像也要加载

	//选加载  做好兜底,万一 Redis真的崩了,你要保护你的mysql
	//我数据库限流

	//选不加载, 用户体验差一点
}
func (r *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

func (r *CachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "", //确实有邮箱
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "", //确实有手机号
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}

func (ur *CachedUserRepository) Update(ctx context.Context, u domain.User) error {
	err := ur.dao.UpdateNonZeroFields(ctx, ur.domainToEntity(u))
	if err != nil {
		return err
	}
	return ur.cache.Delete(ctx, u.Id) // Changed from ur to ur
}
