package v1

import "github.com/gin-gonic/gin"

func Pong(c *gin.Context) {
	c.JSON(200, struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		Code: 0,
		Msg:  "Pong",
	})
}
