package events

import (
	"Book_Exp/webook/feed/domain"
	"Book_Exp/webook/feed/service"
	"Book_Exp/webook/pkg/logger"
	"Book_Exp/webook/pkg/saramax"
	"context"
	"time"

	"github.com/IBM/sarama"
)

const topicFeedEvent = "feed_event"

type FeedEvent struct {
	Type     string
	Metadata map[string]string
}

type FeedEventConsumer struct {
	svc    service.FeedService
	client sarama.Client
	l      logger.LoggerV1
}

func NewFeedEventConsumer(svc service.FeedService, client sarama.Client, l logger.LoggerV1) *FeedEventConsumer {
	return &FeedEventConsumer{
		svc:    svc,
		client: client,
		l:      l,
	}

}
func (c *FeedEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(topicFeedEvent, c.client)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(), []string{topicFeedEvent},
			saramax.NewHandler[FeedEvent](c.l, c.Consume))
		if err != nil {
			c.l.Error("推出了消费循环异常", logger.Error(err))
		}
	}()
	return err

}
func (c *FeedEventConsumer) Consume(sg *sarama.ConsumerMessage,
	evt FeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return c.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: evt.Type,
		Ext:  evt.Metadata,
	})
}
