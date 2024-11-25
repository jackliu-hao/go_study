package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello go",
		})
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		c.JSON(200, gin.H{
			"username": username,
			"password": password,
		})
	})
	//go func() {
	//	engine := gin.Default()
	//	engine.Run(":8081")
	//}()
	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
