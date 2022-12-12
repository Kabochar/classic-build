package main

import (
	"gin/v1/server"
)

func main() {
	engine := server.NewRouter()
	engine.Run() // listen and serve on 0.0.0.0:8080
}
