package service

import (
	"context"
	"errors"
	"log"

	"github.com/jinzhu/copier"

	"gin/v1/models"
	"gin/v1/pkg/typex"
)

type CreateArticleService struct {
	Title   string         `json:"title" form:"title" `
	Alias   typex.StrSlice `json:"alias" form:"alias"`
	Author  string         `json:"author" form:"author"`
	Content string         `json:"content" form:"content"`
	Tags    typex.StrSlice `json:"tags" form:"tags"`
}

func (svc *CreateArticleService) CreateArticle(ctx context.Context) interface{} {
	var article models.Article
	copier.Copy(&article, svc)

	if err := models.DB.WithContext(ctx).Create(&article).Error; err != nil {
		log.Println("create record error: ", err)
		return errors.New("create failed")
	}
	
	return article
}
