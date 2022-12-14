package service

import (
	"context"
	"log"

	"gorm.io/gorm"

	"gin/v1/models"
	"gin/v1/pkg/typex"
)

type ListArticleService struct {
	ID      []int64        `json:"ID" form:"ID"`
	Title   string         `json:"title" form:"title"`
	Alias   typex.StrSlice `json:"alias" form:"alias"`
	Author  string         `json:"author" form:"author"`
	Content string         `json:"content" form:"content"`
	Tags    typex.StrSlice `json:"tags" form:"tags"`
}

func (svc *ListArticleService) ListArticle(ctx context.Context) ([]*models.Article, int64, error) {
	var (
		articles []*models.Article
		count    int64
		err      error
	)

	tx := models.DB.WithContext(ctx).Model(&models.Article{})
	err = buildQuerySQL(svc, tx).Find(&articles).Error
	if err != nil {
		log.Println("find record err", err)
		return nil, count, err
	}

	err = tx.Count(&count).Error
	if err != nil {
		log.Println("count record err", err)
		return nil, count, err
	}
	return articles, count, nil
}

// gorm 多条件查询
// https://blog.csdn.net/weixin_45604257/article/details/106063381
func buildQuerySQL(svc *ListArticleService, tx *gorm.DB) *gorm.DB {
	if len(svc.ID) > 0 {
		tx = tx.Where("id in ?", svc.ID)
	}
	if len(svc.Title) > 0 {
		tx = tx.Where("title like ?", svc.Title)
	}
	if len(svc.Alias) > 0 {
		tx = tx.Where("title like ?", svc.Alias)
	}
	if len(svc.Author) > 0 {
		tx = tx.Where("author = ?", svc.Author)
	}
	if len(svc.Content) > 0 {
		tx = tx.Where("content like ?", svc.Content)
	}
	if len(svc.Tags) > 0 {
		tx = tx.Where("content like ?", svc.Tags)
	}
	return tx
}
