package dao

import (
	"context"

	"gorm.io/gorm"
)

type FeedPushEventDAO interface {
	//创建推送事件
	CreatePushEvents(ctx context.Context, events []FeedPushEvent) error
	GetPushEvents(ctx context.Context, uid int64, timestamp, limit int64) ([]FeedPushEvent, error)
	GetPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp, limit int64) ([]FeedPushEvent, error)
}

// FeedPushEvent 写扩散,推模型,收件箱
// 这个表理论上只插入,不更新,也不删除的 但是可以归档
type FeedPushEvent struct {
	//Id      int64
	//Uid     int64
	//Content string //大的json串
	//Ctime   time.Time
	Id  int64 `gorm:"primaryKey,autoIncrement"`
	UID int64 `gorm:"column:uid;type:int(11);not null;"`
	//用来标记是什么类类型的事件 ,决定了content 怎样解读
	Type string `gorm:"column:type;type:varchar(255);comment:类型"`
	//大的json字符串
	Content string `gorm:"column:content;type:text;"`
	// 发生时间
	Ctime int64 `gorm:"column:ctime;comment:发生时间"`
}
type feedPushEventDAO struct {
	db *gorm.DB
}

func (f *feedPushEventDAO) GetPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp, limit int64) ([]FeedPushEvent, error) {
	var events []FeedPushEvent
	err := f.db.WithContext(ctx).
		Where("type =?", typ).
		Where("uid =?	", uid).
		Where("ctime < ? ", timestamp).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func NewFeedPushEventDAO(db *gorm.DB) FeedPushEventDAO {
	return &feedPushEventDAO{
		db: db,
	}
}

func (f *feedPushEventDAO) CreatePushEvents(ctx context.Context, events []FeedPushEvent) error {
	return f.db.WithContext(ctx).Create(events).Error
}

func (f *feedPushEventDAO) GetPushEvents(ctx context.Context, uid int64, timestamp, limit int64) ([]FeedPushEvent, error) {
	var events []FeedPushEvent
	err := f.db.WithContext(ctx).Where(
		"uid=?", uid).
		Where("ctime<?", timestamp).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}
