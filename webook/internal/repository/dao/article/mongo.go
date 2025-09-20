package article

import (
	"context"
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
}
type MongoDBDAO struct {
	client    *mongo.Client     //连接
	col       *mongo.Collection //代表制作库
	liveCol   *mongo.Collection //代表线上库
	databases *mongo.Database   //代表webook
	node      *snowflake.Node
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) MongoArticleDAO {
	return &MongoDBDAO{
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		node:    node,
	}
}

// 创建
func (m *MongoDBDAO) Insert(ctx context.Context, art Article) (int64, error) {
	id := m.node.Generate().Int64()
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)
	if err != nil {
		return 0, err
	}
	return id, err
}

// 更新
func (m *MongoDBDAO) UpdateById(ctx context.Context, art Article) error {
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	update := bson.D{bson.E{Key: "$set", Value: bson.M{
		"title":   art.Title,
		"content": art.Content,
		"utime":   art.Utime,
		"status":  art.Status,
	}}}
	upRes, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	//这边校验了author_id 是不是正确的id
	if upRes.ModifiedCount == 0 {
		return errors.New("更新失败")
	}
	return nil

}

// 同步
func (m *MongoDBDAO) Sync(ctx context.Context, art Article) (int64, error) {
	//没有办法引入事务 先保存制作库
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return id, err
	}
	art.Id = id
	//操作线上库
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	now := time.Now().UnixMilli()
	art.Utime = now
	update := bson.E{Key: "$set", Value: PublishedArticle(art)}
	sert := bson.E{Key: "$setOnInsert", Value: bson.D{bson.E{Key: "ctime", Value: now}}}
	_, err = m.liveCol.UpdateOne(ctx, filter, bson.D{update, sert}, options.Update().SetUpsert(true))
	return id, err

}
