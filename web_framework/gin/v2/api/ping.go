package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, time.Now().Format("2006-01-02 15:03:04.000")+": pong.....")
}
