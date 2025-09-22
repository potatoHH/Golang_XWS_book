package cache

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInteractiveCache interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrReadCntPresent(ctx context.Context, biz string, bizId int64) error
}

var (
	//go:embed lua/interactive_incr_cnt.lua
	luaIncrCnt string
)

const (
	fieldReadCnt    = "read_cnt"
	fieldCollectCnt = "collect_cnt"
	fieldLikeCnt    = "like_cnt"
)

type RedisInteractive struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewRedisInteractiveCache(client redis.Cmdable) RedisInteractiveCache {
	return &RedisInteractive{
		client: client,
	}
}
func (r *RedisInteractive) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	panic("")
}

func (r *RedisInteractive) IncrReadCntPresent(ctx context.Context, biz string, bizId int64) error {
	//拿到的结果,可能自增成功了,可能不需要自增(key不存在)
	return r.client.Eval(ctx, luaIncrCnt, []string{r.Key(biz, bizId)},
		fieldReadCnt, 1).Err()
	//你哟啊不要返回一个error 表达 key 不存在? 所以不需要
	//if err != nil {
	//	return err
	//}
	//if res == 0 {
	//	return errors.New("缓存中的key 不存在")
	//}
	//return nil

}

func (r *RedisInteractive) Key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
