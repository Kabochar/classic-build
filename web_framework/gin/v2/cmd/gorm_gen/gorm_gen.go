package main

import (
	"log"

	"gin/v2/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "dal",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	db, err := gorm.Open(sqlite.Open(".demo.db"))
	if err != nil {
		log.Println("gorm open error", err)
	}
	g.UseDB(db) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(model.User{})

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	g.ApplyInterface(func(Querier) {}, model.User{}, model.Video{})

	// Generate the code
	g.Execute()

	// auto generate table
	db.AutoMigrate(model.User{}, model.Video{})
}
