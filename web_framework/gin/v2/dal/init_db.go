package dal

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 在中间件中初始化mysql链接
func NewDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             1 * time.Microsecond, // Slow SQL threshold
			LogLevel:                  logger.Info,          // Log level
			IgnoreRecordNotFoundError: true,                 // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open(".demo.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Printf("connect mysql client ERROR %v\n", err)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("get mysql client ERROR %v\n", err)
		return
	}
	// 设置连接池
	// 空闲
	sqlDB.SetMaxIdleConns(50)
	// 打开
	sqlDB.SetMaxOpenConns(100)
	// 超时
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	// 配置全局db client
	SetDefault(db)
}
