package service

import (
	service2 "Book_Exp/webook/interactive/service"
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository"
	"context"

	"math"
	"time"

	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
)

type RankingService interface {
	TopN(ctx context.Context) error
	topN(ctx context.Context) ([]domain.Article, error)
}
type BatchrankingService struct {
	artSvc    ArticleService
	intrSvc   service2.InteractiveService
	repo      repository.RankingRepository
	batchSize int
	n         int
	//scoreFunc 不能返回负数
	scoreFunc func(t time.Time, LikeCnt int64) float64
}

func NewBatchRankingService(artSvc ArticleService, intrSvc service2.InteractiveService) RankingService {
	return &BatchrankingService{
		intrSvc:   intrSvc,
		artSvc:    artSvc,
		batchSize: 100,
		n:         100,
		scoreFunc: func(t time.Time, LikeCnt int64) float64 {
			sec := time.Since(t).Seconds()
			return float64(LikeCnt-1) / math.Pow(float64(sec+2), 1.5)
		},
	}
}
func (svc *BatchrankingService) TopN(ctx context.Context) error {
	arts, err := svc.topN(ctx)
	if err != nil {
		return err
	}
	//放入redis里面
	return svc.repo.ReplaceTopN(ctx, arts)
}
func (svc *BatchrankingService) topN(ctx context.Context) ([]domain.Article, error) {
	//TODO 先拿一批数据
	offset := 0
	type Score struct {
		score float64
		art   domain.Article
	}
	//这里可以用非并发安全数据据
	topN := queue.NewConcurrentPriorityQueue[Score](svc.n, func(src Score, dst Score) int {
		if src.score > dst.score {
			return 1
		} else if src.score == dst.score {
			return 0
		} else {
			return -1
		}
	})
	for {
		//TODO 这里拿了一批
		arts, err := svc.artSvc.ListPub(ctx, nil, offset, svc.batchSize)
		if err != nil {
			return nil, err
		}
		ids := slice.Map[domain.Article, int64](arts, func(idx int, src domain.Article) int64 {
			return src.Id
		})
		//TODO  要去找对应的点赞数据
		intrs, err := svc.intrSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}
		//TODO 要去找到对应的点赞数据
		svc.intrSvc.GetByIds(ctx, "article", []int64{1, 2, 3})
		//TODO 合并计算 score

		//TODO 排序 sort
		for _, art := range arts {
			intr, ok := intrs[art.Id]
			if !ok {
				continue
			}
			score := svc.scoreFunc(art.Utime, intr.LikeCnt)
			//要考虑这个socore在不在前一百名 拿到热度最低的
			if err == queue.ErrOutOfCapacity {
				//这种写法要求topN已经满了
				val, _ := topN.Dequeue()
				if val.score < score {
					topN.Enqueue(Score{
						art:   art,
						score: score,
					})
				}
			}

		}
		//TODO 一批已经处理完成,要不要进行下一批
		now := time.Now()
		if len(arts) < svc.batchSize || now.Sub(arts[0].Utime).Hours() > 7*24 {
			//这一批都没取够,
			break
		}
		//TODO 更新offest
		offset = offset + len(arts)
	}

	//TODO最后得出结果
	res := make([]domain.Article, 0, svc.n)
	for i := svc.n; i > 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			break
		}
		res[i] = val.art
	}
	return res, nil

}
