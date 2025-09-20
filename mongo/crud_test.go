package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
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
	assert.NoError(t, err)
	//InsertedId 文档Id 就是mongodb中的_id 字段
	fmt.Printf("插入了Id:%s", res.InsertedID)

	//TODO bson 来构造查找条件
	filter := bson.D{bson.E{Key: "id", Value: 134}}
	findRes := col.FindOne(ctx, filter)
	var art Article
	err = findRes.Decode(&art) //解析结果
	assert.NoError(t, err)
	fmt.Printf("%v:", art)
	findRes = col.FindOne(ctx, Article{Id: 134})
	//用findRes.Err()==mongo.ErrNodocuments 来判断有没有查找数据  查找不到是没有忽略零值
	if findRes.Err() == mongo.ErrNoDocuments {
		fmt.Println("没有找到")
	} else {
		fmt.Printf("find %#v:", art)
	}
	//TODO Dson 更新文档
	sets := bson.D{bson.E{Key: "$set", Value: bson.E{Key: "title", Value: "新的标题"}}} //只更新这里提到的字段
	updateRes, err := col.UpdateOne(ctx, filter, sets)
	assert.NoError(t, err)
	fmt.Println("affected", updateRes.ModifiedCount)
	//updateMany 其他字段被更新为0值  实践用这种写法比较好
	updateManyRes, err := col.UpdateMany(ctx, filter, bson.D{bson.E{Key: "$set", Value: Article{
		Title: "新的标题2",
	}}})
	assert.NoError(t, err)
	fmt.Printf("更新字段 %v", updateManyRes)
	//TODO 删除文档
	delRes, err := col.DeleteOne(ctx, filter)
	assert.NoError(t, err)
	fmt.Println("删除了", delRes.DeletedCount)

	//TODO OR 查找
	//or 中用了 A 包裹D E
	or := bson.A{bson.D{bson.E{Key: "id", Value: 134}},
		bson.D{bson.E{Key: "id", Value: 123}},
	}
	orRes, err := col.Find(ctx, bson.D{bson.E{Key: "$or", Value: or}})
	assert.NoError(t, err)
	var artsor []Article
	err = orRes.All(ctx, &artsor)
	assert.NoError(t, err)
	//TODO And 查找
	and := bson.A{bson.D{bson.E{Key: "id", Value: 134}},
		bson.D{bson.E{Key: "title", Value: "新的标题"}},
	}
	andRes, err := col.Find(ctx, bson.D{bson.E{Key: "$and", Value: and}})
	assert.NoError(t, err)
	var artsand []Article
	err = andRes.All(ctx, &artsand)
	assert.NoError(t, err)
	//TODO  In 查找
	inRes, err := col.Find(ctx, bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: []any{123, 134}}}}})
	assert.NoError(t, err)
	var artsin []Article
	err = inRes.All(ctx, &artsin)
	assert.NoError(t, err)
	//TODO 查询特定的字段 projection
	proRes, err := col.Find(ctx, bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: []int{123, 134}}}}},
		options.Find().SetProjection(bson.D{
			bson.E{Key: "id", Value: 1},
		}),
	)
	assert.NoError(t, err)
	var artpro []Article
	err = proRes.All(ctx, &artpro)
	assert.NoError(t, err)

	////TODO mongo创建索引
	//idxRes, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
	//	Keys:    bson.D{bson.E{Key: "id", Value: 1}},
	//	Options: options.Index().SetUnique(true), //设置唯一索引
	//})
	//assert.NoError(t, err)
	//fmt.Println("索引创建成功", idxRes)
}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId string `bson:"author_id,omitempty"`
	Status   int64  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
