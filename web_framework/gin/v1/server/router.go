package server

import (
	"github.com/gin-gonic/gin"

	v1 "gin/v1/api/v1"
)

// NewRouter add new gin engine
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", v1.Pong)

	v1Group := r.Group("/api/v1")
	{
		v1Group.GET("/article/:id", v1.GetArticle)
		v1Group.GET("/articles", v1.ListArticle)
		v1Group.POST("/article", v1.CreateArticle)
		v1Group.PUT("/article/:id", v1.UpdateArticle)
		v1Group.DELETE("/article/:id", v1.DeleteArticle)
	}
	return r
}
