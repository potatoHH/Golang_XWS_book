package service

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("用户名或者密码不对")
)

type UserServiceV1 interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type UserService struct { // 用户服务
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserServiceV1 { // 构造函数
	return &UserService{
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
func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.FindOrCreate(ctx, phone)
	//判断有没有这个用户  //快路径
	if err != repository.ErrUserNotFound {
		return u, err
	}
	//TODO 在系统资源不足,触发降级之后,不执行慢路径
	//if ctx.Value("降级") == "true" {
	//	return domain.User{}, errors.New("系统资源不足 降级")
	//}
	//你要明确,没有这个用户  慢路径
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	//这个u没有id 因为这里会遇到主从延迟的问题
	return svc.repo.FindByPhone(ctx, phone)

}

func (svc *UserService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	// 写法1
	// 这种是简单的写法，依赖与 Web 层保证没有敏感数据被修改
	// 也就是说，你的基本假设是前端传过来的数据就是不会修改 Email，Phone 之类的信息的。
	//return svc.repo.Update(ctx, user)

	// 写法2
	// 这种是复杂写法，依赖于 repository 中更新会忽略 0 值
	// 这个转换的意义在于，你在 failover 层面上维护住了什么是敏感字段这个语义
	user.Email = ""
	user.Phone = ""
	user.Password = ""
	return svc.repo.Update(ctx, user)
}

// 信息
func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	//在缓存中取 在系统内部一般都是使用id 来取
	//有些人的系统比较复杂,有一个GUId(global unique id )  全局唯一
	return svc.repo.FindById(ctx, id)
}
