package article

import (
	"Book_Exp/webook/interactive/repository"
	"Book_Exp/webook/pkg/logger"

	"github.com/IBM/sarama"
)

type HisotryReadEventConsumer struct {
	l      logger.LoggerV1
	client sarama.Client
	repo   repository.InteractiveRepository
}

func NewHisotryReadEventConsumer(l logger.LoggerV1, repo repository.InteractiveRepository, client sarama.Client) *HisotryReadEventConsumer {
	return &HisotryReadEventConsumer{
		l:      l,
		client: client,
		repo:   repo,
	}
}
