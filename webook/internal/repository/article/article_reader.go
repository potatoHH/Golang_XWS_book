package article

import (
	"Book_Exp/webook/internal/domain"
	"context"
)

type ArticleReaderRepository interface {
	//有就更新,没有就新建 即upsert的语句
	Save(ctx context.Context, art domain.Article) (int64, error)
}
