package metric

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type MiddlewareBuilder struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	InstanceID string
}

func (m *MiddlewareBuilder) Build() gin.HandlerFunc {
	//pattern 是命中路由
	labels := []string{"method", "pattern", "status"}
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name,
		ConstLabels: prometheus.Labels{
			"instance_id": m.InstanceID,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
			1:     0.0001,
		},
	}, labels)
	prometheus.MustRegister(summary)
	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			pattern := ctx.FullPath()
			if pattern == "" {
				pattern = "unknow"
			}
			summary.WithLabelValues(ctx.Request.Method, pattern, strconv.Itoa(ctx.Writer.Status())).Observe(float64(duration.Milliseconds()))
			summary.WithLabelValues(ctx.Request.Method, "", "").Observe(float64(duration.Milliseconds()))
		}() //你最终会执行到业务里面
		ctx.Next()
	}

}
