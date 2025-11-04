package article

import (
	"Book_Exp/webook/internal/repository"
	"Book_Exp/webook/pkg/logger"
	"Book_Exp/webook/pkg/saramax"
	"context"
	"time"

	"github.com/IBM/sarama"
)

type InteractiveReadEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewInteractiveReadEventConsumer(l logger.LoggerV1, repo repository.InteractiveRepository, client sarama.Client) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}

func (r *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"artcile_read"},
			saramax.NewHandler[ReadEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常",
				logger.Error(err))
		}
	}()
	return err

}

func (r *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()
	return r.repo.IncrReadCnt(ctx, "article", t.Aid)
}
