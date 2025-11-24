package dao

import (
	"Book_Exp/webook/internal/repository/dao"
	"context"
	"errors"

	"github.com/ecodeclub/ekit/syncx/atomicx"
)

var errUnknownPattern = errors.New("未知的双写 pattern")

const (
	patternSrcOnly  = "SRC_ONLY"
	patternSrcFirst = "SRC_FIRST"
	patternDstOnly  = "DST_ONLY"
	patternDstFirst = "DST_FIRST"
)

type DoubleWirteDAO struct {
	src     InteractiveDAO
	dst     InteractiveDAO
	pattern *atomicx.Value[string]
}

func NewDoubleWirteDAO(src InteractiveDAO, dst InteractiveDAO) *DoubleWirteDAO {
	return &DoubleWirteDAO{src: src, dst: dst,
		pattern: atomicx.NewValueOf(patternDstFirst)}
}

func (d *DoubleWirteDAO) UpdatePattern(pattern string) {
	d.pattern.Store(pattern)
}
func (d *DoubleWirteDAO) DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	switch d.pattern.Load() {
	case patternSrcOnly:
		return d.src.DeleteLikeInfo(ctx, biz, id, uid)
	case patternSrcFirst:
		err := d.src.DeleteLikeInfo(ctx, biz, id, uid)
		if err != nil {
			return err
		}
		//这是一个问题, src 成功了,但是dst失败了怎么?,等修复
		return d.dst.DeleteLikeInfo(ctx, biz, id, uid)
	case patternDstOnly:
		return d.dst.DeleteLikeInfo(ctx, biz, id, uid)
	case patternDstFirst:
		err := d.dst.DeleteLikeInfo(ctx, biz, id, uid)
		if err != nil {
			return err
		}
		return d.dst.DeleteLikeInfo(ctx, biz, id, uid)
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWirteDAO) InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	switch d.pattern.Load() {
	case "SRC_ONLY":
		return d.src.DeleteLikeInfo(ctx, biz, id, uid)
	case "SRC_FIRST":
		err := d.src.DeleteLikeInfo(ctx, biz, id, uid)
		if err != nil {
			return err
		}
		//这是一个问题, src 成功了,但是dst失败了怎么?,等修复
		err = d.dst.DeleteLikeInfo(ctx, biz, id, uid)
		if err != nil {
			//记录日志
		}
		return nil
	case "DST_ONLY":
		return d.dst.DeleteLikeInfo(ctx, biz, id, uid)
	case "DST_FIRST":
		err := d.dst.DeleteLikeInfo(ctx, biz, id, uid)
		if err != nil {
			return err
		}
		err = d.src.DeleteLikeInfo(ctx, biz, id, uid)
		if err != nil {
			//记录日志
		}
		return nil
	default:
		return errors.New("未知的双写模式")
	}
}

func (d *DoubleWirteDAO) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	switch d.pattern.Load() {
	case patternSrcOnly:
		return d.src.IncrReadCnt(ctx, biz, id)
	case patternSrcFirst:
		err := d.src.IncrReadCnt(ctx, biz, id)
		if err != nil {
			return err
		}
		err = d.dst.IncrReadCnt(ctx, biz, id)
		if err != nil {
			//记录日志
		}
		return nil
	case patternDstOnly:
		return d.dst.IncrReadCnt(ctx, biz, id)
	case patternDstFirst:
		err := d.dst.IncrReadCnt(ctx, biz, id)
		if err != nil {
			return err
		}
		err = d.src.IncrReadCnt(ctx, biz, id)
		if err != nil {
			//记录日志
		}
		return nil
	default:
		return errors.New("未知的双写模式")
	}
}

func (d *DoubleWirteDAO) InsertCollectionBiz(ctx context.Context, biz dao.UserCollectionBiz) error {
	//TODO implement me
	panic("implement me")
}
func (d *DoubleWirteDAO) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	switch d.pattern.Load() {
	case patternSrcOnly, patternSrcFirst:
		return d.src.Get(ctx, biz, id)
	case patternDstOnly, patternDstFirst:
		return d.dst.Get(ctx, biz, id)
	default:
		return Interactive{}, errUnknownPattern
	}

}

func (d *DoubleWirteDAO) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (dao.UserLikeBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWirteDAO) GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (dao.UserCollectionBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWirteDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, bizId []int64) error {
	//TODO implement me
	panic("implement me")
}
