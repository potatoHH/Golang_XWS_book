package ratelimit

import (
	"context"
)

type Limiter interface {
	//TODO 有没有触发限流,key就是限流对象  bool是否限流  err 限流器本身有没有错误
	Limit(ctx context.Context, key string) (bool, error)
}

//type RedisSlidingWindowLimiter struct {
//	cmd      redis.Cmdable
//	interval time.Duration // 滑动窗口的间隔
//	rate     int
//}
//
//func NewRedsisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
//	return &RedisSlidingWindowLimiter{
//		cmd:      cmd,
//		interval: interval,
//		rate:     rate,
//	}
//}
//
//func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool ,error) {
//	return r.cmd.Eval(ctx, luaScript, []string{key}),
//		r.interval.Milliseconds(),
//		r.rate, time.Now().UnixMilli())
//}
