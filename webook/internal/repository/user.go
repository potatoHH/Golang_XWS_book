package repository

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository/cache"
	"Book_Exp/webook/internal/repository/dao"
	"context"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct { // 用户仓库
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {

		}
	}()
	return u, nil

	//这里怎么版?err =io.EOF
	//要不要去数据库加载?
	//看起来我不应该加载
	//看起来我好像也要加载

	//选加载  做好兜底,万一 Redis真的崩了,你要保护你的mysql
	//我数据库限流

	//选不加载, 用户体验差一点

}
