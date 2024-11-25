package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jikeshijian_go/webbook/internal/repository"
	"jikeshijian_go/webbook/internal/repository/dao"
	"jikeshijian_go/webbook/internal/service"
	"jikeshijian_go/webbook/internal/web"
	"jikeshijian_go/webbook/internal/web/middleware"
	"strings"
	"time"
)

func main() {

	// 初始化db
	db := InitDb()

	// 初始化UserHandler
	uh := InitUserHandler(db)

	// 初始化gin Engine
	server := InitWebServer()

	// 注册路由
	uh.RegisterRoutes(server)
	server.Run(":8081")

}

func InitUserHandler(db *gorm.DB) *web.UserHandler {
	// 初始化dao层
	ud := dao.NewUserDAO(db)

	// 初始化respo
	repo := repository.NewUserRepository(ud)

	// 初始化service
	svc := service.NewUserService(repo)

	// 创建uerHandler
	uh := web.NewUserHandler(svc)

	return uh
}

func InitWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true, // 携带cookie

		AllowHeaders: []string{"Content-Type"},
		//AllowHeaders: []string{"content-type"},
		//AllowMethods: []string{"POST"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	// 初始化session的store
	// 基于cookie的store
	//store := cookie.NewStore([]byte("secret_key_1234567890"))
	// 基于内存实现的store
	//store := memstore.NewStore([]byte("0776f450dd575004ba7c69930c579cae"),
	//	[]byte("0776f450dd575004ba7c69930c579cae"))
	// 基于redis实现的store
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("0776f450dd575004ba7c69930c579cae"),
		[]byte("0776f450dd575004ba7c69930c579cae"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))
	// 注册登录的middleware
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePath("/users/signup").
		IgnorePath("/users/login").
		Build())
	return server
}

func InitDb() *gorm.DB {
	// 初始化db
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/webook"))
	if err != nil {
		panic(err)
	}
	// 初始化db
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
