package service

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户名或者密码不对")
)

type UserService struct { // 用户服务
	repo *repository.UserRepository
}

func NewService(repo *repository.UserRepository) *UserService { // 构造函数
	return &UserService{ // 返回一个用户服务
		repo: repo,
	}

}

// 注册
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost) // 密码加密
	if err != nil {
		return err
	}
	u.Password = string(hash) // 将密码加密后的结果赋给用户

	return svc.repo.Create(ctx, u)
}

// 登录
func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	//先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		//打日志
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

// 信息
func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	//在缓存中取
	u, err := svc.repo.FindById(ctx, id)
	return u, err
}
