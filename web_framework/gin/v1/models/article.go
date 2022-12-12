package models

import (
	"gorm.io/gorm"

	"gin/v1/pkg/typex"
)

type Article struct {
	gorm.Model
	ID      int `gorm:"primaryKey;autoIncrement"`
	Title   string
	Alias   typex.StrSlice
	Author  string
	Content string `gorm:"size:1000"`
	Tags    typex.StrSlice
}
