package service

import (
	"context"
	"log"

	"gin/v1/models"
)

type GetArticleService struct {
}

func (svc *GetArticleService) GetArticle(ctx context.Context, id int64) interface{} {
	var article models.Article

	if err := models.DB.WithContext(ctx).First(&article, id).Error; err != nil {
		log.Println("find record err", err)
		return err
	}

	return article
}
