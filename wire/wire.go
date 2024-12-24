//go:build wireinject

// 让 wire 来注入这里的代码

package wire

import (
	"github.com/google/wire"
	"jikeshijian_go/wire/repository"
	"jikeshijian_go/wire/repository/dao"
)

func InitUserRepository() *repository.UserRepository {
	// 这个方法里面传入各个组件的初始化方法，wire会自动根据依赖关系进行组装
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)

	return &repository.UserRepository{}
}
