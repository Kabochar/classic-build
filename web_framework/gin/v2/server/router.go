package server

import (
	"github.com/gin-gonic/gin"

	"gin/v2/api"
)

// NewRouter return server all api define
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", api.Pong)
	r.POST("/api/v2/video/create", api.CreateVideo)
	return r
}
