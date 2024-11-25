package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"jikeshijian_go/webbook/internal/domain"
	"jikeshijian_go/webbook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

// Create 创建
func (r *UserRepository) Create(ctx context.Context, user domain.User) error {

	// 1. 密码加密

	// 2. 插入到数据库

	// 3. 返回结果
	return r.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})

}

func NewUserRepository(dao *dao.UserDAO) UserRepository {

	return UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

	// 先从cache中查找

	// 再从dao中查找
	user, err := r.dao.FindById(ctx, id)
	return user, err

	// 找到了在写到cache

}

func (r *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	// 需要将PO转成BO ,返回给service层
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) UpdateById(ctx *gin.Context, user domain.User) error {
	err := r.dao.UpdateById(ctx, dao.User{
		Id:       user.Id,
		NickName: user.NickName,
		AboutMe:  user.AboutMe,
		Birthday: user.Birthday,
	})
	return err
}
