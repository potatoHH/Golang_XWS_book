package ioc

import (
	"Book_Exp/webook/interactive/events"
	"Book_Exp/webook/pkg/saramax"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

//func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
//	res, err := sarama.NewSyncProducerFromClient(client)
//	if err != nil {
//		panic(err)
//	}
//	return res
//}

// TODO NewConsumers 面临的问题依旧是所有的Consumer 在这里注册一下
func NewConsumers(c1 *events.InteractiveReadEventConsumer) []saramax.Consumer {
	return []saramax.Consumer{
		c1,
	}
}
