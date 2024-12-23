package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱或手机号冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func (dao *UserDAO) Insert(ctx context.Context, user User) error {
	// 拿到时间戳的毫秒数
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&user).Error

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		//  mysql唯一索引错误码
		const uniqueConflictsErr = 1062
		if mysqlErr.Number == uniqueConflictsErr {
			// 邮箱冲突
			return ErrUserDuplicate
		}
	}

	return err

}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {

	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) UpdateById(ctx context.Context, user User) error {

	err := dao.db.WithContext(ctx).Model(&user).Updates(user).Error
	return err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	// select * from `user` where id = ? limit 1

	err := dao.db.WithContext(ctx).First(&user, "id = ?", id).Error

	//err := dao.db.WithContext(ctx).Model(User{Id: id}).First(&user).Error
	// 将int64转成time类型
	return user, err
}

func (dao *UserDAO) FindByPhone(ctx *gin.Context, phone string) (User, error) {
	var user User
	// select * from `user` where id = ? limit 1

	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error

	return user, err
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// User : entity 、 model 、 po
// 对应DDD中的entity
type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	// 昵称
	NickName string
	// 个人简介
	AboutMe string
	// 生日
	Birthday string
	// 手机号
	// 唯一索引，允许存在多个空值
	// 但是不能允许多个 ""
	Phone sql.NullString `gorm:"unique"`

	// 创建时间，毫秒数
	CreatedAt int64
	// 更新时间，毫秒数
	UpdatedAt int64
}

// Encrypt 加密方法
//func (u *User) Encrypt() {
//
//}
