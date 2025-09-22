package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
}

type GormInteractiveDAO struct {
	db *gorm.DB
}

// IncrReadCnt 是一个插入或者更新语义
func (I *GormInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	now := time.Now().UnixMilli()
	return I.db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt": gorm.Expr("read_cnt + 1"),
			"utime":    now,
		}),
	}).Create(&Interactive{
		Biz:     biz,
		BizId:   id,
		Ctime:   now,
		Utime:   now,
		ReadCnt: 1,
	}).Error

}

func NewGormInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GormInteractiveDAO{
		db: db,
	}
}

type Interactive struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gomr:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCnt    int64
	CollectCnt int64
	LikeCnt    int64
	Ctime      int64
	Utime      int64
}
