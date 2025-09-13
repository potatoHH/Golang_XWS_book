package article

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository/dao/article"
	"context"

	"gorm.io/gorm"
)

// repository 还是用来操作缓存和dao
//事务概念应该在dao中定义

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error) //创建数据
	Update(ctx context.Context, art domain.Article) error              //更新数据
	SyncV1(ctx context.Context, art domain.Article) (int64, error)     //存储并同步数据
}
type CacheArticleRepostior struct {
	dao article.ArticleDAO
	//耦合了DAO 操作的东西,
	//正常情况下,如果你要在reposiotry 层面上操作事务
	//那么这就是能利用db开始事务之后,创建基于事务的DAO
	//或者,直接去掉DAO 这一层,在respository的实现中,直接操作db
	db *gorm.DB

	//v1 操作两个dao
	authorDao article.ArticleAuthorDAO
	readerDao article.ArticleReaderDAO
}

func NewArticleRepostior(dao article.ArticleDAO) ArticleRepository {
	return &CacheArticleRepostior{
		dao: dao,
	}
}
func (c *CacheArticleRepostior) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.toEntity(art))

}

// 创建
func (c *CacheArticleRepostior) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, article.Article{
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
	})
}

// 更新
func (c *CacheArticleRepostior) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, article.Article{
		Id:       art.Id,
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
	})
}

// TODO 在Syncv2 尝试再repository层上解决事务问题 ,确保制作库和线上库同时成功,或者同时失败
func (c *CacheArticleRepostior) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	//开启事务
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()
	//利用tx来构造 dao
	author := article.NewGORMArticleDAO(tx)
	reader := article.NewGORMArticleReaderDAO(tx)
	artn := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = author.UpdateById(ctx, artn)
	} else {
		id, err = author.Insert(ctx, artn)
	}
	if err != nil {
		//执行有问题,需要回滚
		return id, err
	}
	err = reader.UpsertV2(ctx, article.PublishArticle{Article: artn})
	tx.Commit()
	return id, err

}

func (c *CacheArticleRepostior) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	artn := c.toEntity(art)
	//TODO  应该先保存到制作库,再保存的线上库
	if art.Id > 0 {
		err = c.authorDao.UpdateById(ctx, artn)

	} else {
		id, err = c.authorDao.Insert(ctx, artn)
	}
	if err != nil {
		return id, err
	}
	//TODO 操作线上库,保存数据,同步过来,考虑到,此时线上库可能有或没有,你要有一个upsert方法 INSERT OR Update , 如果数据库有就更新,没有就插入
	err = c.readerDao.UpSert(ctx, artn)
	return id, err
}

func (c *CacheArticleRepostior) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
	}

}
