package job

import (
	"Book_Exp/webook/internal/service"
	"context"
	"time"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)

}
