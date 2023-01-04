package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	ID    int64  `json:"ID"`
	Title string `json:"title"`
	Alias string `json:"alias"`
	Tags  string `json:"tags"`
	Desc  string `json:"desc"`
}
