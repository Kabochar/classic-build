package server

import "gin/v1/models"

func InitServer() {
	models.NewDatabase("")
}
