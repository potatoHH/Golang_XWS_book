package validator

import (
	"Book_Exp/webook/migrator"
	"Book_Exp/webook/migrator/events"
	"Book_Exp/webook/pkg/logger"
	"context"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// t 必须实现了migrator.Entity 接口
type Validator[T migrator.Entity] struct {
	l      logger.LoggerV1
	base   *gorm.DB
	target *gorm.DB
	//这边需要告知是以SRC为准,还是以DST为准
	//修复数据需要知道
	direction string
	batchSize int
	p         events.Producer
}

func (v *Validator[T]) validate(ctx context.Context)error {
	// 创建一个 errgroup 来并发执行两个验证任务
	// errgroup 可以方便地管理多个 goroutine，并收集它们的错误
	eg := errgroup.Group{}
	
	// 第一个 goroutine：验证从 base 到 target 的数据一致性
	// 检查 base 中的每条记录在 target 中是否存在且一致
	eg.Go(func() error {
		v.ValidateBaseToTarget(ctx)
		return nil
	})
	
	// 第二个 goroutine：验证从 target 到 base 的数据一致性
	// 检查 target 中是否存在 base 中已删除的数据
	eg.Go(func() error {
		v.ValidateTargetToBase(ctx)
		return nil
	})
	
	// 等待所有 goroutine 完成并返回可能的错误
	// 如果任何一个 goroutine 返回错误，Wait() 会返回第一个非 nil 的错误
	return eg.Wait()

}

// Validate 验证者可以通过ctx来验证退出
func (v *Validator[T]) ValidateBaseToTarget(ctx context.Context) {
	//进来就更新,比较好控制
	offset := -1
	for {
		dbCtx, cancle := context.WithTimeout(ctx, time.Second*10)
		offset++
		var src T
		//例如.Order("id Desc"),这样插入数据 offset 将会不准确
		err := v.base.WithContext(dbCtx).Offset(offset).First(&src).Order("id").Error
		cancle()
		switch err {
		case nil:
			//查找到了数据,继续去target中查找
			var dst T
			err1 := v.target.Where("id=?", src.ID()).First(&dst).Error
			switch err1 {
			case nil:
				//找到了开始比较,这是利用 反射来比较
				if !src.CompareTo(dst) {
					//这时候要上报给kafka,告诉他数据不一致
					v.nofity(ctx, src.ID(), events.InconsistentEventyTypeNEQ)

				}
				//var srcAny any = src
				//if c1, ok := srcAny.(interface {
				//	//有没有自定义的比较逻辑
				//	CompareTo(c2 migrator.Entity) bool
				//}); ok {
				//	//有的话就用他的
				//	if c1.CompareTo(dst) {
				//		return
				//	} else {
				//		//没有的话就用反射
				//		if !reflect.DeepEqual(src, dst) {
				//			return
				//		}
				//	}
				//}
			case gorm.ErrRecordNotFound:
				//意味着 target里面少了数据
				v.nofity(ctx, src.ID(), events.InconsistentEventTypeTargetMissing)
				return
			default:
				//1 我认为,大概率数据是一致的,我记录一下日志,下体条
				v.l.Error("校验数据,查询target失败", logger.Error(err1))
				continue
				//2.我认为出入保险起见,我应该报数据不一致,试着去修改一下, 如果真的不一致,就修改,如果数据一直,就多修改一次

			}
		case gorm.ErrRecordNotFound:
			//比完了,没有数据 全量校验结束
			return
		default:
			v.l.Error("校验数据,查询base 失败", logger.Error(err))
			continue
		}

	}

}

func (v *Validator[T]) nofity(ctx context.Context, id int64, typ string) {
	ctx, cancle := context.WithTimeout(ctx, time.Second*10)
	err := v.p.ProducerInconsistentEvent(ctx, events.InconsistentEvent{
		ID:        id,
		Direction: v.direction,
		Type:      typ,
	})
	cancle()
	if err != nil {
		//失败了,记录日志,警告,手动去修改
		v.l.Error("上报数据不一致事件失败", logger.Error(err))
	}
}
func (v *Validator[T]) nofityBaseMissing(ctx context.Context, ids []int64) {
	for _, id := range ids {
		v.nofity(ctx, id, events.InconsistentEventTypeBaseMissing)

	}

}

func (v *Validator[T]) ValidateTargetToBase(ctx context.Context) {
	//先找target 在找base, 找出base 已被删除的数据
	offset := -v.batchSize
	for {
		dbCtx, cancle := context.WithTimeout(ctx, time.Second*10)
		offset := offset + v.batchSize
		var dstTs []T
		//例如.Order("id Desc"),这样插入数据 offset 将会不准确
		err := v.base.WithContext(dbCtx).Offset(offset).Limit(v.batchSize).Find(&dstTs).Order("id").Error
		cancle()
		if len(dstTs) == 0 {
			return
		}
		switch err {
		case nil:
			ids := slice.Map(dstTs, func(idx int, t T) int64 {
				return t.ID()

			})
			var srcTs []T
			err := v.base.Where("id IN ?", ids).Find(&srcTs).Error
			switch err {
			case gorm.ErrRecordNotFound:
				v.nofityBaseMissing(ctx, ids)
			case nil:
				srcIds := slice.Map(srcTs, func(idx int, t T) int64 {
					return t.ID()
				})
				//计算差值 ,也即是src里面没有的数据
				diff := slice.DiffSet(ids, srcIds)
				v.nofityBaseMissing(ctx, diff)
			default:
				v.l.Error("校验数据,查询base 失败", logger.Error(err))
				continue
			}
		case gorm.ErrRecordNotFound:
		default:
			v.l.Error("校验数据,查询base 失败", logger.Error(err))
			continue
		}
		if len(dstTs) < v.batchSize {
			return
		}
	}

}
