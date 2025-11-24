package fixer

import (
	"Book_Exp/webook/migrator"
	"Book_Exp/webook/migrator/events"
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Finxer[T migrator.Entity] struct {
	base    *gorm.DB
	target  *gorm.DB
	columns []string
}

// 一了百了,直接覆盖  event只是一个触发器, 不依赖里面的具体内容,ID 必须不可变
func (f *Finxer[T]) Fix(ctx context.Context, evt events.InconsistentEvent) error {
	var t T
	err := f.base.WithContext(ctx).Where("id =?", evt.ID).First(&t).Error
	switch err {
	case nil:
		//base 里面有数据 ,就更新
		return f.target.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(f.columns)}).
			Where("id=?", evt.ID).First(&t).Error
	case gorm.ErrRecordNotFound:
		//base 没有
		return f.target.WithContext(ctx).Where("id=?", evt.ID).Delete(&t).Error

	default:
		return err

	}
}

// base 和 target 在校验时候的数据,到你修复的时候就变了
func (f *Finxer[T]) FixV1(ctx context.Context, evt events.InconsistentEvent) error {
	switch evt.Type {
	case events.InconsistentEventTypeTargetMissing, events.InconsistentEventyTypeNEQ:
		//这边要插入
		var t T
		err := f.base.WithContext(ctx).Clauses(
			clause.OnConflict{
				//更新全部列
				DoUpdates: clause.AssignmentColumns(f.columns),
			}).
			Where("id=?", evt.ID).First(&t).Error
		switch err {
		case nil:
			//就在你插入的时候,双写的程序,也插入了,就会冲突
			return f.target.Create(&t).Error
		case gorm.ErrRecordNotFound:
			return f.target.Where("id=?", evt.ID).Delete(new(T)).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		return f.target.Where("id=?", evt.ID).Delete(new(T)).Error
	default:
		return errors.New("未知的不一致类型")
	}
}

func (f *Finxer[T]) FixV2(ctx context.Context, evt events.InconsistentEvent) error {
	switch evt.Type {
	case events.InconsistentEventTypeTargetMissing:
		//这边要插入
		var t T
		err := f.base.WithContext(ctx).Where("id=?", evt.ID).First(&t).Error
		switch err {
		case nil:
			return f.target.Create(&t).Error
		case gorm.ErrRecordNotFound:
			return nil
		default:
			return err
		}

	case events.InconsistentEventyTypeNEQ:
		//这边要删除
		var t T
		err := f.target.WithContext(ctx).Where("id=?", evt.ID).First(&t).Error
		switch err {
		case nil:
			return f.target.Updates(&t).Error
		case gorm.ErrRecordNotFound:
			return f.target.Delete(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		return f.target.Where("id=?", evt.ID).Delete(new(T)).Error
	default:
		return errors.New("未知的不一致类型")
	}
}
