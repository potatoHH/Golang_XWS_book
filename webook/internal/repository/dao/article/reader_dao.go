package article

import (
	"context"

	"gorm.io/gorm"
)

type ArticleReaderDAO interface {
	UpSert(ctx context.Context, art Article) error
	UpsertV2(ctx context.Context, art PublishArticle) error
}

// PublishArticle 这个代表线上库
type PublishArticle struct {
	Article
}

type GORMArticleReaderDAO struct {
	db *gorm.DB
}

func NewGORMArticleReaderDAO(db *gorm.DB) ArticleReaderDAO {
	panic("")
}

func (dao *GORMArticleReaderDAO) UpSert(ctx context.Context, art Article) error {
	panic("")

}
