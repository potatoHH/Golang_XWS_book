package dao

import "gorm.io/gorm"

//使用gorm的建表功能

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{}) // 自动迁移

}
