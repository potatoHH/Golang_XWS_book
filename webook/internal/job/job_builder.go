package job

import (
	"Book_Exp/webook/pkg/logger"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
)

type CronJobBuilder struct {
	l logger.LoggerV1
	p *prometheus.SummaryVec
}

func NewCronJobBuilder(l logger.LoggerV1, opt prometheus.SummaryOpts) *CronJobBuilder {
	p := prometheus.NewSummaryVec(opt, []string{"name"})
	prometheus.MustRegister(p)
	return &CronJobBuilder{
		l: l,
		p: p,
	}

}

func (r *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobFuncAdapter(func() error {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Milliseconds()
			r.p.WithLabelValues(name).Observe(float64(duration))
		}()
		err := job.Run()
		if err != nil {
			r.l.Error("任务执行失败", logger.Error(err), logger.String("job", job.Name()))
		}
		return nil
	})

}

type cronJobFuncAdapter func() error

func (c cronJobFuncAdapter) Run() {
_:
	c()

}
