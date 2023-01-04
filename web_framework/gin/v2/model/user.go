package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       string `json:"ID"`
	Name     string `json:"name"`
	NickName string `json:"nickName"`
	Mobile   string `json:"mobile"`
	Address  string `json:"address"`
	IP       string `json:"ip"`
}
