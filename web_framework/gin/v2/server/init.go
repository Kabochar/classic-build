package server

import "gin/v2/dal"

// Run init hole server
func Run() {
	dal.NewDatabase()

	r := NewRouter()
	r.Run()
}
