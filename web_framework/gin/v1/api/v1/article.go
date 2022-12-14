package v1

import (
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"gin/v1/service"
)

func CreateArticle(c *gin.Context) {
	var svc service.CreateArticleService
	ctx := context.Background()

	err := c.ShouldBind(&svc)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200,
		svc.CreateArticle(ctx),
	)
}

func UpdateArticle(c *gin.Context) {
	var svc service.UpdateArticleService
	ctx := context.Background()

	err := c.ShouldBind(&svc)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200,
		svc.UpdateArticle(ctx, id),
	)
}

func GetArticle(c *gin.Context) {
	var svc service.GetArticleService
	ctx := context.Background()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200,
		svc.GetArticle(ctx, id),
	)
}

func ListArticle(c *gin.Context) {
	var svc service.ListArticleService
	ctx := context.Background()

	err := c.ShouldBind(&svc)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	result, count, err := svc.ListArticle(ctx)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	log.Println("get record: ", count)
	c.JSON(200,
		result,
	)
}

func DeleteArticle(c *gin.Context) {
	var svc service.DeleteArticleService
	ctx := context.Background()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200,
		svc.DeleteArticle(ctx, id),
	)
}
