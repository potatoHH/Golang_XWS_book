package dao

import (
	"Book_Exp/webook/internal/repository/dao/article"

	"gorm.io/gorm"
)

//使用gorm的建表功能

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&article.Article{},
		&article.PublishArticle{},
	) // 自动迁移

}
