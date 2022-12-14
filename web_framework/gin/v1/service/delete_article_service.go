package service

import (
	"context"
	"log"

	"gin/v1/models"
)

type DeleteArticleService struct {
}

func (DeleteArticleService) DeleteArticle(ctx context.Context, id int64) interface{} {
	var article models.Article

	tx := models.DB.WithContext(ctx)
	if err := tx.First(&article, id).Error; err != nil {
		log.Printf("find record err: {%+v}", err)
		return err
	}

	if err := tx.Delete(&article, id).Error; err != nil {
		log.Printf("delete record : {%+v}", err)
		return err
	}

	return nil
}
