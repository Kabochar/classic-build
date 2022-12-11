package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "pong",
			"demo":    "中文字符串处理",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
