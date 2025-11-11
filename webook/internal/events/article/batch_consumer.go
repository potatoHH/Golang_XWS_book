package article

import (
	"Book_Exp/webook/interactive/repository"
	"Book_Exp/webook/pkg/logger"
	"Book_Exp/webook/pkg/saramax"
	"context"
	"time"

	"github.com/IBM/sarama"
)

type InteractiveReadEventBatchConsumer struct {
	l      logger.LoggerV1
	client sarama.Client
	repo   repository.InteractiveRepository
}

func NewInteractiveReadEventBatchConsumer(l logger.LoggerV1, repo repository.InteractiveRepository, client sarama.Client) *InteractiveReadEventBatchConsumer {
	return &InteractiveReadEventBatchConsumer{
		l:      l,
		client: client,
		repo:   repo,
	}
}

func (r *InteractiveReadEventBatchConsumer) StartBatch() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"artcile_read"},
			saramax.NewBatchHandler[ReadEvent](r.l, r.ConsumeBatch))
		if err != nil {
			r.l.Error("退出了消费循环异常",
				logger.Error(err))
		}
	}()
	return err

}
func (r *InteractiveReadEventBatchConsumer) ConsumeBatch(msg []*sarama.ConsumerMessage, ts []ReadEvent) error {
	bizs := make([]string, 0, len(msg))
	ids := make([]int64, 0, len(msg))
	for _, evt := range ts {
		ids = append(ids, evt.Aid)
		bizs = append(bizs, "article")
	}
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	err := r.repo.BatchReadCnt(ctx, bizs, ids)
	if err != nil {
		r.l.Error("批量更新阅读数失败",
			logger.Field{Key: "ids", Value: "ids"},
			logger.Error(err))
	}
	cancle()
	return nil
}
