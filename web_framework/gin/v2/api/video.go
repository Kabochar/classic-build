package api

import (
	"context"
	"log"
	"net/http"

	"gin/v2/service"

	"github.com/gin-gonic/gin"
)

func CreateVideo(c *gin.Context) {
	var svc service.CreateVideoService
	if err := c.ShouldBind(&svc); err != nil {
		log.Println("bind error", err)
		c.JSON(http.StatusOK, gin.H{
			"err": err,
		})
		return
	}

	ctx := context.Background()
	result, err := svc.CreateVideo(ctx)
	if err != nil {
		log.Println("svc.CreateVideo error", err)
		c.JSON(http.StatusOK, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func GetVideo(c *gin.Context) {
	// todo impl
}

func ListVideo(c *gin.Context) {
	// todo impl
}

func UpdateVideo(c *gin.Context) {
	// todo impl
}

func DeleteVideo(c *gin.Context) {
	// todo impl
}
