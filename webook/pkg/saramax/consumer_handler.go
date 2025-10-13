package saramax

import (
	"Book_Exp/webook/pkg/logger"
	"encoding/json"

	"github.com/IBM/sarama"
)

type Handler[T any] struct {
	l  logger.LoggerV1
	fn func(msg *sarama.ConsumerMessage, t T) error
}

func NewHandler[T any](l logger.LoggerV1, fn func(msg *sarama.ConsumerMessage, t T) error) *Handler[T] {
	return &Handler[T]{
		l:  l,
		fn: fn,
	}
}

func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		//TODO  反序列化失败,说明数据有问题,重试没有必要
		if err != nil {
			h.l.Error("序列化失败", logger.Error(err),
				logger.String("topic", msg.Topic),
				logger.Int64("partition", int64(msg.Partition)),
				logger.Int64("offset", msg.Offset))
			continue
		}
		err = h.fn(msg, t)
		//这里执行重试
		if err != nil {
			h.l.Error("序列化失败", logger.Error(err),
				logger.String("topic", msg.Topic),
				logger.Int64("partition", int64(msg.Partition)),
				logger.Int64("offset", msg.Offset))

		}
	}
	return nil
}
