package connpool

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
)

var errUnknownPattern = errors.New("未知的双写错误")

const (
	patternSrcOnly  = "SRC_ONLY"
	patternDstOnly  = "DST_ONLY"
	patternSrcFirst = "SRC_FIRST"
	patternDstFirst = "DST_FIRST"
)

type DoubleWirtePool struct {
	src     gorm.ConnPool
	dst     gorm.ConnPool
	pattern *atomicx.Value[string]
}

func (d *DoubleWirtePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	//TODO sql.Stmt 是一个结构体,你没有办法返回一个代表双写的 stmt
	panic("双写模式下不支持")

}

func (d *DoubleWirtePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	//增删改,或者非查询语句
	switch d.pattern.Load() {
	case patternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case patternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		//这是一个问题, src 成功了,但是dst失败了怎么?,等修复
		return d.dst.ExecContext(ctx, query, args...)
	case patternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	case patternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		_, err = d.src.ExecContext(ctx, query, args...)
		if err != nil {
			//dst 写错了,不认为是错误
		}
		return res, err
	default:
		panic("未知的双写粗模式")
		return nil, errUnknownPattern
	}
}

func (d *DoubleWirtePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern.Load() {
	case patternSrcOnly, patternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		panic("未知的双写粗模式")
		return nil, errUnknownPattern
	}
}

func (d *DoubleWirtePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern.Load() {
	case patternSrcOnly, patternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		panic("未知的双写粗模式")
	}
}

type DoubleWirtePoolTx struct {
	src *sql.Tx
	dst *sql.Tx
	//**从事务一致性角度考虑**：
	//- 事务开始时应该确定使用哪种模式
	//- 事务执行过程中，模式不应该改变
	//- 因此事务对象应该持有事务开始时的模式快照，而不是共享引用
	pattern string
}

func (d *DoubleWirtePoolTx) Commit() error {
	switch d.pattern {
	case patternSrcOnly:
		return d.src.Commit()
	case patternSrcFirst:
		//源库的数据提交失败了,目标库需不需要提交
		err := d.src.Commit()
		if err != nil {
			return err
		}
		if d.dst != nil {
			err = d.dst.Commit()
			if err != nil {

			}
		}
		return nil
	case patternDstOnly:
		return d.dst.Commit()
	case patternDstFirst:
		err := d.dst.Commit()
		if err != nil {
			return err
		}
		if d.src != nil {
			err = d.src.Commit()
			if err != nil {
			}
		}
		return nil
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWirtePoolTx) Rollback() error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWirtePool) BeginTx(ctx context.Context, opt *sql.TxOptions) (gorm.ConnPool, error) {
	var pattern = d.pattern.Load()
	switch pattern {
	case patternSrcOnly:
		tx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opt)
		return &DoubleWirtePoolTx{src: tx}, err
	case patternSrcFirst:
		srcTx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opt)
		if err != nil {
			return nil, err
		}
		dstTx, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opt)
		if err != nil {
			//记录日志;这里不做处理
		}
		return &DoubleWirtePoolTx{src: srcTx, dst: dstTx, pattern: pattern}, nil

	case patternDstOnly:
		tx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opt)
		return &DoubleWirtePoolTx{src: tx, pattern: pattern}, err
	case patternDstFirst:
		dstTx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opt)
		if err != nil {
			return nil, err
		}
		srcTx, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opt)
		if err != nil {
			//记录日志;这里不做处理
		}
		return &DoubleWirtePoolTx{
			src:     srcTx,
			dst:     dstTx,
			pattern: pattern,
		}, nil
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWirtePoolTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	//TODO sql.Stmt 是一个结构体,你没有办法返回一个代表双写的 stmt
	panic("双写模式下不支持")
}

func (d *DoubleWirtePoolTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	//增删改,或者非查询语句
	switch d.pattern {
	case patternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case patternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		if d.dst == nil {
			return res, nil

		}
		//这是一个问题, src 成功了,但是dst失败了怎么?,等修复
		return d.dst.ExecContext(ctx, query, args...)
	case patternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	case patternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		if d.src == nil {
			return res, nil
		}
		_, err = d.src.ExecContext(ctx, query, args...)
		if err != nil {
			//dst 写错了,不认为是错误
		}
		return res, err
	default:
		panic("未知的双写粗模式")
		return nil, errUnknownPattern
	}

}

func (d *DoubleWirtePoolTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	//查询语句
	switch d.pattern {
	case patternSrcOnly, patternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		panic("未知的双写粗模式")
		return nil, errUnknownPattern
	}
}

func (d *DoubleWirtePoolTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	//查询语句
	switch d.pattern {
	case patternSrcOnly, patternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		panic("未知的双写粗模式")
	}
}
