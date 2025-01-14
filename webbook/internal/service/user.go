package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"jikeshijian_go/webbook/internal/domain"
	"jikeshijian_go/webbook/internal/repository"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账号或密码错误")
	ErrUserNotFind           = repository.ErrUserNotFound
)

type UserService interface {
	SingUp(ctx context.Context, user domain.User) error
	Login(ctx context.Context, u domain.User) (domain.User, error)
	Edit(ctx context.Context, user domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(c context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(c context.Context, info domain.WechatInfo) (domain.User, error)
}

type UserServiceV1 struct {
	userRepository repository.UserRepository
}

func (svc *UserServiceV1) FindOrCreateByWechat(c context.Context, info domain.WechatInfo) (domain.User, error) {
	// 快路径
	user, err := svc.userRepository.FindByWechat(c, info.OpenId)
	// 需要判断这个用户是否存在，如果不存在需要创建
	if !errors.Is(err, ErrUserNotFind) {
		// 其他非预期错误
		return user, err
	} else {
		// 确实是没找到有这个用户，需要创建一个用户
		user = domain.User{
			WechatInfo: info,
		}
		// 慢路径
		err := svc.userRepository.Create(c, user)
		if err != nil && !errors.Is(err, ErrUserDuplicate) {
			return domain.User{}, err
		}
		// 再次找一下，并且返回
		// 可能存在主从延迟的问题
		return svc.userRepository.FindByWechat(c, info.OpenId)
	}

}

func NewUserServiceV1(userRepository repository.UserRepository) UserService {
	return &UserServiceV1{
		userRepository: userRepository,
	}
}

func (svc *UserServiceV1) SingUp(ctx context.Context, user domain.User) error {
	// 1. 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	// 2. 保存到数据库
	return svc.userRepository.Create(ctx, user)
}

func (svc *UserServiceV1) Login(ctx context.Context, u domain.User) (domain.User, error) {
	// 查询用户
	user, err := svc.userRepository.FindByEmail(ctx, u.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.User{}, ErrUserNotFind
		}
		return domain.User{}, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		// 账号或密码错误
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func (svc *UserServiceV1) Edit(ctx context.Context, user domain.User) error {
	err := svc.userRepository.UpdateById(ctx, user)
	return err
}

func (svc *UserServiceV1) Profile(ctx context.Context, id int64) (domain.User, error) {
	user, err := svc.userRepository.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (svc *UserServiceV1) FindOrCreate(c context.Context, phone string) (domain.User, error) {

	// 快路径
	user, err := svc.userRepository.FindByPhone(c, phone)
	if !errors.Is(err, ErrUserNotFind) {
		// 其他错误
		return user, err
	} else {
		// 需要判断这个用户是否存在，如果不存在需要创建
		user = domain.User{
			Phone: phone,
		}
		// 慢路径
		err := svc.userRepository.Create(c, user)
		if err != nil && !errors.Is(err, ErrUserDuplicate) {
			return domain.User{}, err
		}
		// 再次找一下，并且返回
		// 可能存在主从延迟的问题
		return svc.userRepository.FindByPhone(c, phone)
	}
}
