//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"jikeshijian_go/webbook/internal/repository"
	"jikeshijian_go/webbook/internal/repository/cache"
	"jikeshijian_go/webbook/internal/repository/dao"
	"jikeshijian_go/webbook/internal/service"
	"jikeshijian_go/webbook/internal/web"
	"jikeshijian_go/webbook/ioc"
	"time"
)

func InitWebServerIOC() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		// DAO 部分
		dao.NewGormUserDAO,

		// cache 部分
		InitCacheTime,
		cache.NewRedisCodeCache, cache.NewRedisUserCache,

		// repository 部分
		repository.NewUserRepositoryWithCache,
		repository.NewCodeRepositoryWithCache,

		// Service 部分
		ioc.InitSMSService,
		service.NewUserServiceV1,
		InitCodeServiceTpl,
		service.NewCodeServiceWith6Num,

		// handler 部分
		web.NewUserHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}

func InitCacheTime() time.Duration {
	return time.Duration(5 * time.Minute)
}

func InitCodeServiceTpl() string {
	return "testId"
}
