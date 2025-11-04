package cron

import (
	"log"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

type JobFunc struct{}

func (j JobFunc) Run() {
	log.Println("运行了")
}

func TestCronExpr(t *testing.T) {
	expr := cron.New(cron.WithSeconds())
	//这个是任务的标识符
	//@Every是便利的语法
	//expr.AddJob("@every 1s", JobFunc{})
	expr.AddFunc("@every 1s", func() {
		log.Println("开始长任务")
		time.Sleep(15 * time.Second)
		log.Println("我也运行了")
	})
	//调度运行
	expr.Start()
	time.Sleep(10 * time.Second)
	//发出停止信号,不会调度新的任务

	ctx := expr.Stop()
	log.Println("发出停止信号")
	<-ctx.Done()
	log.Println("彻底结束")
}
