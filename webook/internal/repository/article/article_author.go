package article

import (
	"Book_Exp/webook/internal/domain"
	"context"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}
