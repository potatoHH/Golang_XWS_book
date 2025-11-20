package article

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository"
	"Book_Exp/webook/internal/repository/cache"
	"Book_Exp/webook/internal/repository/dao/article"
	"Book_Exp/webook/pkg/logger"
	"context"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
)

// repository 还是用来操作缓存和dao
//事务概念应该在dao中定义

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error) //创建数据
	Update(ctx context.Context, art domain.Article) error              //更新数据
	SyncV1(ctx context.Context, art domain.Article) (int64, error)     //存储并同步数据
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset int, limt int) ([]domain.Article, error)

	//收藏
}
type CacheArticleRepostiory struct {
	dao      article.ArticleDAO
	userRepo repository.UserRepository
	//耦合了DAO 操作的东西,
	//正常情况下,如果你要在reposiotry 层面上操作事务
	//那么这就是能利用db开始事务之后,创建基于事务的DAO
	//或者,直接去掉DAO 这一层,在respository的实现中,直接操作db
	db    *gorm.DB
	cache cache.ArticleCache
	//v1 操作两个dao
	authorDao article.ArticleAuthorDAO
	readerDao article.ArticleReaderDAO
	l         logger.LoggerV1
}

func (repo *CacheArticleRepostiory) ListPub(ctx context.Context, start time.Time, offset int, limt int) ([]domain.Article, error) {
	_, err := repo.dao.ListPub(ctx, start, offset, limt)
	if err != nil {
		return nil, err
	}
	return slice.Map[article.PublishedArticle, domain.Article](nil, func(idx int, src article.PublishedArticle) domain.Article {
		// 偷懒写法
		return repo.ToDomain(article.Article(src))
	}), nil
}

func (repo *CacheArticleRepostiory) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	//读取线上库数据,如果你的content 被你放过去了 oss上,你就要让前端去读content字段
	art, err := repo.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	//这边要组装user,适合单体应用
	usr, err := repo.userRepo.FindById(ctx, art.AuthorId)
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id:   usr.Id,
			Name: usr.Nickname,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	return res, nil

}
func (repo *CacheArticleRepostiory) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	data, err := repo.dao.GetByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return repo.ToDomain(data), nil
}

func NewCacheArticleRepostiory(dao article.ArticleDAO, l logger.LoggerV1) ArticleRepository {
	return &CacheArticleRepostiory{
		dao: dao,
		l:   l,
	}
}

// 列表
func (repo *CacheArticleRepostiory) List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error) {
	//TODO 在这个地方,集成你的复杂的缓存方案
	if offset == 0 && limit == 100 {
		data, err := repo.cache.GetFirstPage(ctx, uid)
		if err == nil {
			go func() {
				repo.PreCache(ctx, data)
			}()
			return data, err
		}
	}
	res, err := repo.dao.GetByAuthor(ctx, uid, limit, offset)
	if err != nil {
		return nil, err
	}
	data := slice.Map[article.Article, domain.Article](res, func(idx int, src article.Article) domain.Article {
		return repo.ToDomain(src)
	})
	//TODO 回写缓存的时候,你要考虑是Set 还是Del,如果你认为,比如说创作者不太可能有很高并发,你就直接用set ,如果要是有很高并发就用del
	//回写缓存的时候 ,可以同步,也可以异步
	go func() {
		err := repo.cache.SetFirstPage(ctx, uid, data)
		repo.l.Error("回写缓存失败", logger.Error(err))
		repo.PreCache(ctx, data)
	}()
	return data, err

}

func (repo *CacheArticleRepostiory) SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, id, author, status.ToUnit8())
}

func NewArticleRepostior(dao article.ArticleDAO) ArticleRepository {
	return &CacheArticleRepostiory{
		dao: dao,
	}
}
func (repo *CacheArticleRepostiory) Sync(ctx context.Context, art domain.Article) (int64, error) {
	//TODO 清空缓存
	id, err := repo.dao.Sync(ctx, repo.toEntity(art))
	if err == nil {
		//提前缓存好线上库数据
		repo.cache.DelFirstPage(ctx, art.Author.Id)
		repo.cache.SetPub(ctx, art)
		if err != nil {
			repo.l.Error("提前设置缓存失败", logger.Int64("author", art.Author.Id), logger.Error(err))
		}
	}
	return id, err

}

// 创建
func (repo *CacheArticleRepostiory) Create(ctx context.Context, art domain.Article) (int64, error) {
	//清空缓存
	defer func() {
		repo.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return repo.dao.Insert(ctx, article.Article{
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
		Status:   art.Status.ToUnit8(),
	})
}

// 更新
func (repo *CacheArticleRepostiory) Update(ctx context.Context, art domain.Article) error {
	return repo.dao.UpdateById(ctx, article.Article{
		Id:       art.Id,
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
		Status:   art.Status.ToUnit8(),
	})
}

// TODO 在Syncv2 尝试再repository层上解决事务问题 ,确保制作库和线上库同时成功,或者同时失败
func (repo *CacheArticleRepostiory) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	//开启事务
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()
	//利用tx来构造 dao
	author := article.NewGORMArticleDAO(tx)
	reader := article.NewGORMArticleReaderDAO(tx)
	artn := repo.toEntity(art)
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

func (repo *CacheArticleRepostiory) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	artn := repo.toEntity(art)
	//TODO  应该先保存到制作库,再保存的线上库
	if art.Id > 0 {
		err = repo.authorDao.UpdateById(ctx, artn)

	} else {
		id, err = repo.authorDao.Insert(ctx, artn)
	}
	if err != nil {
		return id, err
	}
	//TODO 操作线上库,保存数据,同步过来,考虑到,此时线上库可能有或没有,你要有一个upsert方法 INSERT OR Update , 如果数据库有就更新,没有就插入
	err = repo.readerDao.UpSert(ctx, artn)
	return id, err
}

func (repo *CacheArticleRepostiory) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
		Status:   art.Status.ToUnit8(),
	}

}
func (repo *CacheArticleRepostiory) ToDomain(art article.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
}

// TODO 提前预加载缓存
func (repo *CacheArticleRepostiory) PreCache(ctx context.Context, data []domain.Article) {
	const contentSizeThreshold = 1024 * 1024
	if len(data) > 0 && len(data[0].Content) < contentSizeThreshold {
		//你也可以记录日志
		if err := repo.cache.Set(ctx, data[0].Id); err != nil {
			repo.l.Error("提前准备缓存失败", logger.Error(err))
		}
	}
}
