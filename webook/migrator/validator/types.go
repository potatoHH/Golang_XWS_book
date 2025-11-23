package validator

import (
	"Book_Exp/webook/migrator"
	"Book_Exp/webook/pkg/logger"
	"context"

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
}

// Validate 验证者可以通过ctx来验证退出
func (v *Validator[T]) Validate(ctx context.Context) {
	//进来就更新,比较好控制
	offset := -1
	for {
		offset++
		var src T
		//例如.Order("id Desc"),这样插入数据 offset 将会不准确
		err := v.base.Offset(offset).First(&src).Order("id").Error
		switch err {
		case nil:
			//查找到了数据,继续去target中查找
			var dst T
			v.target.Where("id=?", src.ID()).First(&dst)
		case gorm.ErrRecordNotFound:
			//比完了,没有数据 全量校验结束
			return
		default:
			v.l.Error("校验数据,查询base 失败", logger.Error(err))
			continue
		}

	}

}
