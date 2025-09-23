package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
}

type GormInteractiveDAO struct {
	db *gorm.DB
}

func (dao *GormInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId int64, uid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("biz=? AND bizId=? AND Uid =?", biz, bizId, uid).Updates(map[string]any{
			"utime":  now,
			"status": 0,
		}).Error
		if err != nil {
			return err
		}
		return tx.Where("biz =? AND bizId =?", biz, bizId).Updates(map[string]any{
			"utime":    now,
			"like_cnt": gorm.Expr("llike_cnt - 1"),
		}).Error
	})

}

// 点赞
func (dao *GormInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO 一把梭, 同时记录点赞,以及更新点赞计数,首先你需要一张表来记录,谁点给什么资源点了赞
	//TODO 创建一个事务
	now := time.Now().UnixMilli()
	return dao.db.Model(&UserLikeBiz{}).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"status": 1,
				"utime":  now,
			}),
		}).Create(&UserLikeBiz{
			Uid:    uid,
			Ctime:  now,
			Utime:  now,
			Biz:    biz,
			BizId:  id,
			Status: 1,
		}).Error
		if err != nil {
			return err
		}

		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt": gorm.Expr("llike_cnt + 1"),
				"utime":    now,
			}),
		}).Create(&Interactive{
			Biz:     biz,
			BizId:   id,
			Ctime:   now,
			Utime:   now,
			LikeCnt: 1,
		}).Error
	})
}

// IncrReadCnt 是一个插入或者更新语义
func (dao *GormInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	now := time.Now().UnixMilli()
	return dao.db.Clauses(clause.OnConflict{
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

// 用户交互
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

// 用户点赞的某个东西
type UserLikeBiz struct {
	Id int64 `gorm:"primarykey,autoIncrement"`
	//三个构造一个唯一索引
	BizId int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
	Uid   int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	//依旧是指在DB层面生效的状态
	//1- 有效 0----无效 软删除的用法
	Status uint8
	Ctime  int64
	Utime  int64
}

//用户如何查找前一百条点赞数据
//TODO  Interactive 正常来说,一张主表和其他表的关联关系的表会公用一个DAO ,所以我们就用过一个DAO 来操作 加入说我要查找点赞量前面100
//TODO 实时查找,性能贼查, 高性能,我不要求准确性 面试标准答案用zset,但是不够有特色,烂大街了,你可以考虑 1.定时计算 1.1 定时计算+ 本地缓存
//TODO 2. 优化版zset,定时筛选 zset + 实时 zset计算
