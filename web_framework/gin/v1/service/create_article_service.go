package server

import (
	"fmt"

	"gin/v1/models"

	"github.com/jinzhu/copier"
)

type CreateArticleService struct {
	Title   string   `json:"title" form:"title" `
	Alias   []string `json:"alias" form:"alias"`
	Author  string   `json:"author" form:"author"`
	Content string   `json:"content" form:"content"`
	Tags    []string `json:"tags" form:"tags"`
}

func (svc *CreateArticleService) CreateArticle() {
	var article models.Article
	copier.Copy(article, svc)
	fmt.Printf("%+v", article)
}
