package saramax

import (
	"Book_Exp/webook/pkg/logger"
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

type BatchHandler[T any] struct {
	l             logger.LoggerV1
	fn            func(msgs []*sarama.ConsumerMessage, ts []T) error
	batchSize     int
	batchDuration time.Duration
}

func NewBatchHandler[T any](l logger.LoggerV1, fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{
		l:  l,
		fn: fn,
	}
}
func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgsCh := claim.Messages()
	batchSize := b.batchSize
	for {
		ctx, cancel := context.WithTimeout(context.Background(), b.batchDuration)
		var last *sarama.ConsumerMessage
		msgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				cancel()
				done = true
			case msg, ok := <-msgsCh:
				if !ok {
					cancel()
					//代表消费队列被关闭
					return nil
				}
				last = msg

				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列化失败", logger.Error(err),
						logger.String("msg", string(msg.Value)),
						logger.String("topic", msg.Topic),
						logger.Int64("partition", int64(msg.Partition)),
						logger.Int64("offset", msg.Offset))
					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)

			}
		}
		cancel()
		if len(msgs) == 0 {
			continue
		}
		err := b.fn(msgs, ts)
		if err != nil {
			b.l.Error("处理失败", logger.Error(err))

		}
		if last != nil {
			session.MarkMessage(last, "")
		}
	}
}
