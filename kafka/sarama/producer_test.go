package sarama

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"

	"testing"
)

var adds = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true //同步的produce 一定要设置
	producer, err := sarama.NewSyncProducer(adds, cfg)
	assert.NoError(t, err)
	p, offest, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_group",
		Key:   sarama.StringEncoder("key"),
		Value: sarama.StringEncoder("同步消息"),
		//会在produce 和 consume 中之间传递
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("h1"),
				Value: []byte("h1_value"),
			},
		},
		//只作用于发送过程
		Metadata: "元数据",
	})
	assert.NoError(t, err)
	fmt.Println(p, offest)

}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(adds, cfg)
	assert.NoError(t, err)
	msgCh := producer.Input()
	msgCh <- &sarama.ProducerMessage{
		Topic: "test1",
		Key:   sarama.StringEncoder("key1"),
		Value: sarama.StringEncoder("异步消息"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("h1"),
				Value: []byte("h1_value"),
			},
		},
	}
	//在实践中,一般是开另一个goroutine来处理结果的
	select {
	case err := <-producer.Errors():
		//这边是出错了
		fmt.Println(err)
	case msg := <-producer.Successes():
		//这边是成功了
		fmt.Println(msg)
	}
}
