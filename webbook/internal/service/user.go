package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"jikeshijian_go/webbook/internal/domain"
	"jikeshijian_go/webbook/internal/repository"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账号或密码错误")
	ErrUserNotFind           = repository.ErrUserNotFound
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (svc *UserService) SingUp(ctx context.Context, user domain.User) error {
	// 1. 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	// 2. 保存到数据库
	return svc.userRepository.Create(ctx, user)
}

func (svc *UserService) Login(ctx *gin.Context, u domain.User) (domain.User, error) {
	// 查询用户
	user, err := svc.userRepository.FindByEmail(ctx, u.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return domain.User{}, ErrInvalidUserOrPassword
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

func (svc *UserService) Edit(ctx *gin.Context, user domain.User) error {
	err := svc.userRepository.UpdateById(ctx, user)
	return err
}

func (svc *UserService) Profile(ctx *gin.Context, id int64) (domain.User, error) {
	user, err := svc.userRepository.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (svc *UserService) FindOrCreate(c *gin.Context, phone string) (domain.User, error) {

	// 快路径
	user, err := svc.userRepository.FindByPhone(c, phone)
	if errors.Is(err, ErrUserNotFind) {
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
	} else {
		// 其他错误
		return domain.User{}, err
	}
}
