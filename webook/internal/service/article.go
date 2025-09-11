package service

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository/article"
	"Book_Exp/webook/pkg/logger"
	"context"
	"time"
)

type ArticleService interface {
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
}

type ArticleServiceV1 struct {
	repo article.ArticleRepository

	//TODO V1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
	l      logger.LoggerV1
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &ArticleServiceV1{
		repo: repo,
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
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

func (a *ArticleServiceV1) Publish(ctx context.Context, art domain.Article) (int64, error) {
	panic("123")

	//制作库
	//id, err := a.repo.Create(ctx, art)
	////线上库
	//id, err := a.repo.SyncToLiveDB(ctx, art)

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
