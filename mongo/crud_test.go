package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		//每个命令查询开始之前
		Started: func(c context.Context, cn *event.CommandStartedEvent) {
			fmt.Println(cn.Command)

		},
		//执行成功
		Succeeded: func(c context.Context, cn *event.CommandSucceededEvent) {

		},
		//执行失败
		Failed: func(c context.Context, cn *event.CommandFailedEvent) {

		},
	}
	opts := options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)
	mdb := client.Database("webook")
	col := mdb.Collection("articles")
	res, err := col.InsertOne(ctx, Article{
		Id:      134,
		Title:   "测试",
		Content: "测试内容",
		Status:  1,
		Ctime:   time.Now().UnixMilli(),
		Utime:   time.Now().UnixMilli(),
	})
	if err != nil {
		panic(err)
	}
	//InsertedId 文档Id 就是mongodb中的_id 字段
	fmt.Printf("插入了Id:%s", res.InsertedID)

}

type Article struct {
	Id       int64
	Title    string
	Content  string
	AuthorId string
	Status   int64
	Ctime    int64
	Utime    int64
}
