package service

import (
	"context"
	"log"

	"github.com/jinzhu/copier"

	"gin/v1/models"
	"gin/v1/pkg/typex"
)

type UpdateArticleService struct {
	Title   string         `json:"title" form:"title" `
	Alias   typex.StrSlice `json:"alias" form:"alias"`
	Author  string         `json:"author" form:"author"`
	Content string         `json:"content" form:"content"`
	Tags    typex.StrSlice `json:"tags" form:"tags"`
}

func (svc *UpdateArticleService) UpdateArticle(ctx context.Context, id int64) interface{} {
	var article models.Article

	tx := models.DB.WithContext(ctx)
	if err := tx.First(&article, id).Error; err != nil {
		log.Println("find record err ", err)
		return err
	}

	copier.Copy(&article, svc)
	if err := tx.Updates(&article).Error; err != nil {
		log.Println("update record err ", err)
		return err
	}

	if err := tx.First(&article, id).Error; err != nil {
		log.Println("find record err ", err)
		return err
	}
	return &article
}
