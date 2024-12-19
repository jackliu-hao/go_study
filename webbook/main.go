package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jikeshijian_go/webbook/internal/repository"
	"jikeshijian_go/webbook/internal/repository/cache"
	"jikeshijian_go/webbook/internal/repository/dao"
	"jikeshijian_go/webbook/internal/service"
	"jikeshijian_go/webbook/internal/web"
	"jikeshijian_go/webbook/internal/web/middleware"
	"net/http"
	"strings"
	"time"
)

func main() {

	// 初始化db
	//db := InitDb()
	//
	//// 初始化UserHandler
	//uh := InitUserHandler(db)
	//
	//// 初始化gin Engine
	//server := InitWebServer()
	//
	//// 注册路由
	//uh.RegisterRoutes(server)

	// 用于测试k8s部署
	server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello gin",
			"time":    time.Now().Unix(),
		})
	})
	server.Run(":8081")

}

func InitUserHandler(db *gorm.DB) *web.UserHandler {
	// 初始化dao层
	ud := dao.NewUserDAO(db)

	// todo 初始化cache
	userCache := cache.NewUserCache(nil, time.Minute*5)

	// 初始化resp
	repo := repository.NewUserRepository(ud, userCache)

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
		AllowCredentials: true,                    // 携带cookie
		ExposeHeaders:    []string{"x-jwt-token"}, // 加上这个后，前端才能拿到结果
		AllowHeaders:     []string{"Content-Type", "Authorization"},
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
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("0776f450dd575004ba7c69930c579cae"),
	//	[]byte("0776f450dd575004ba7c69930c579cae"))
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("mysession", store))
	// 注册登录的middleware
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePath("/users/signup").
	//	IgnorePath("/users/login").
	//	Build())
	// 注册jwt的middleware
	server.Use(middleware.NewLoginJwtMiddlewareBuilder().
		IgnorePath("/users/signup").
		IgnorePath("/users/loginjwt").
		Build())
	// 使用限流的工具
	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "",
	//})
	//
	//server.Use(ratelimit.NewBuilder(client, time.Minute, 5).Build())
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
