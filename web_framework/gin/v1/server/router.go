package server

import (
	v1 "gin/v1/api/v1"
	"github.com/gin-gonic/gin"
)

// NewRouter add new gin engine
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", v1.Pong)
	return r
}
