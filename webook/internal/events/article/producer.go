package article

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
)

type Producer interface {
	ProduceReadEvent(ctx context.Context, ent ReadEvent) error
}

type KafkaProducer struct {
	producer  sarama.SyncProducer  //同步的
	aproducer sarama.AsyncProducer //异步的
}

// TODO ProduceReadEvent,如果你有复杂的重试逻辑,就用装饰器 ,你认为你的重试逻辑简单,你就放这里
func (k *KafkaProducer) ProduceReadEvent(ctx context.Context, ent ReadEvent) error {
	data, err := json.Marshal(ent)
	if err != nil {
		return err
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "read-events",
		Value: sarama.ByteEncoder(data),
	})
	return err
}

type ReadEvent struct {
	Uid int64
	Aid int64
}

func NewKafkaProducer(pc sarama.SyncProducer) Producer {
	return &KafkaProducer{
		producer: pc,
	}
}
