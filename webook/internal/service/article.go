package service

import (
	"Book_Exp/webook/internal/domain"
	events "Book_Exp/webook/internal/events/article"
	"Book_Exp/webook/internal/repository/article"
	"Book_Exp/webook/pkg/logger"
	"context"
	"time"
)

//go:generate mockgen -source=./article.go -package=svcmocks -destination=mocks/article.mock.go ArticleService
type ArticleService interface {
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	Save(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
	GetById(ctx context.Context, id, uid int64) (domain.Article, error)
	ListPub(ctx context.Context, offset, limt int) ([]domain.Article, error)
}

type ArticleServiceV1 struct {
	repo article.ArticleRepository
	//TODO V1 依赖两个不同的repository 来解决这种跨表, 或者跨库的问题
	author   article.ArticleAuthorRepository
	reader   article.ArticleReaderRepository
	l        logger.LoggerV1
	producer events.Producer
}

func (a *ArticleServiceV1) ListPub(ctx context.Context, offset, limt int) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (a *ArticleServiceV1) GetById(ctx context.Context, id, uid int64) (domain.Article, error) {
	//另一个选项,在这里组装Author ,调用UserService
	art, err := a.repo.GetByID(ctx, id)
	if err == nil {
		go func() {
			er := a.producer.ProduceReadEvent(ctx, events.ReadEvent{
				//TODO 即便你的消费者要用art里面的数据,让他去查询,你不要在event里面带
				Aid: id,
				Uid: uid,
			})
			if er != nil {
				a.l.Error("发送事件失败",
					logger.Error(er),
				)
			}
		}()
	}
	return art, err
}

func (a *ArticleServiceV1) List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error) {
	return a.repo.List(ctx, uid, limit, offset)
}

func (a *ArticleServiceV1) Withdraw(ctx context.Context, art domain.Article) error {
	return a.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func NewArticleService(repo article.ArticleRepository, producer events.Producer, l logger.LoggerV1) ArticleService {
	return &ArticleServiceV1{
		repo:     repo,
		producer: producer,
		l:        l,
	}
}
func NewArticleServiceV1(author article.ArticleAuthorRepository, reader article.ArticleReaderRepository, l logger.LoggerV1) ArticleService {
	return &ArticleServiceV1{
		author: author,
		reader: reader,
		l:      l,
	}
}

func (a *ArticleServiceV1) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

func (a *ArticleServiceV1) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)
}

func (a *ArticleServiceV1) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	//TODO 确保制作库和线上库的值是相等的
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = a.author.Update(ctx, art)
	} else {
		id, err = a.author.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(i)) //模拟延迟
		id, err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
		a.l.Error("部分失败: 把偶错你数据库到线上失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
	}
	if err != nil {
		a.l.Error("部分失败:重试彻底失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err),
			//接入你的警告系统,手动处理一下
			//走异步,直接保存到本地文件 //走 Canal
		)
	}
	return id, err

}
