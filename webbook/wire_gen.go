// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"jikeshijian_go/webbook/internal/repository"
	"jikeshijian_go/webbook/internal/repository/cache"
	"jikeshijian_go/webbook/internal/repository/dao"
	"jikeshijian_go/webbook/internal/service"
	"jikeshijian_go/webbook/internal/web"
	"jikeshijian_go/webbook/ioc"
	"time"
)

// Injectors from wire.go:

func InitWebServerIOC() *gin.Engine {
	cmdable := ioc.InitRedis()
	v := ioc.InitGinMiddlewares(cmdable)
	db := ioc.InitDB()
	userDao := dao.NewGormUserDAO(db)
	duration := InitCacheTime()
	userCache := cache.NewRedisUserCache(cmdable, duration)
	userRepository := repository.NewUserRepositoryWithCache(userDao, userCache)
	userService := service.NewUserServiceV1(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCodeRepositoryWithCache(codeCache)
	smsService := ioc.InitSMSService()
	string2 := InitCodeServiceTpl()
	codeService := service.NewCodeServiceWith6Num(codeRepository, smsService, string2)
	userHandler := web.NewUserHandler(userService, codeService)
	engine := ioc.InitWebServer(v, userHandler)
	return engine
}

// wire.go:

func InitCacheTime() time.Duration {
	return time.Duration(5 * time.Minute)
}

func InitCodeServiceTpl() string {
	return "testId"
}