package article

import (
	"context"

	"gorm.io/gorm"
)

type ArticleAuthorDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	UpSert(ctx context.Context, art PublishArticle) error
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func (dao *GORMArticleDAO) UpSert(ctx context.Context, art PublishArticle) error {
	//TODO implement me
	panic("implement me")
}

func NewGORMArticleDAO(db *gorm.DB) ArticleAuthorDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	panic("implement me")

}
func (dao *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	panic("implement me")
}
