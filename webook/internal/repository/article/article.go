package article

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository/dao"
	"context"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}
type CacheArticleRepostior struct {
	dao dao.ArticleDAO
}

func NewArticleRepostior(dao dao.ArticleDAO) ArticleRepository {
	return &CacheArticleRepostior{
		dao: dao,
	}
}

// 创建
func (c *CacheArticleRepostior) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		AuthorId: article.Author.Id,
		Content:  article.Content,
		Title:    article.Title,
	})
}

// 更新
func (c *CacheArticleRepostior) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		AuthorId: article.Author.Id,
		Content:  article.Content,
		Title:    article.Title,
	})
}
