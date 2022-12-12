package models

// 自动生成表结构
func migration() {
	_ = DB.AutoMigrate(&Article{})
}
