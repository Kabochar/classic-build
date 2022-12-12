package v1

import (
	"context"

	"github.com/gin-gonic/gin"

	"gin/v1/service"
)

func CreateArticle(c *gin.Context) {
	var svc service.CreateArticleService
	err := c.ShouldBind(&svc)
	if err != nil {
		c.JSON(200, gin.H{})
		return
	}
	c.JSON(200,
		svc.CreateArticle(context.Background()),
	)
}

func UpdateArticle(c *gin.Context) {

}

func ListArticle(c *gin.Context) {

}

func DeleteArticle(c *gin.Context) {

}
