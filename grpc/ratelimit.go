package grpc

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ecodeclub/ekit/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO  限流算法 计算器算法
type CounterLimiter struct {
	cnt       atomic.Int32
	threshold int32
}

func (l *CounterLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		cnt := l.cnt.Add(1)
		defer func() {
			l.cnt.Add(-1)
		}()
		if cnt > l.threshold {
			//这里就是拒绝
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
		}
		return handler(ctx, req)
	}
}

//TODO  固定窗口

type FixedWindowLimiter struct {
	window    time.Duration
	lastTime  time.Time
	cnt       int
	threshold int
	lock      sync.Mutex
}

func (l *FixedWindowLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l.lock.Lock()
		now := time.Now()
		//要换窗口
		if now.Before(l.lastTime.Add(l.window)) {
			l.lastTime = now
			l.cnt = 0
		}
		l.cnt++
		if l.cnt >= l.threshold {
			res, err := handler(ctx, req)
			return res, err
		}
		l.lock.Unlock()
		return nil, status.Errorf(codes.ResourceExhausted, "限流了")

	}
}

//TODO 限流 滑动窗口

// SlidingWindowLimiter 滑动窗口限流器
// 记录每个请求的时间戳，统计滑动窗口内的请求数量
type SlidingWindowLimiter struct {
	window    time.Duration                  // 时间窗口大小
	threshold int                            // 窗口内最大请求数
	requests  []time.Time                    // 请求时间戳队列
	queue     queue.PriorityQueue[time.Time] //
	lock      sync.Mutex
}

func (l *SlidingWindowLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l.lock.Lock()
		now := time.Now()
		windowStart := now.Add(-l.window)
		// 清理窗口之外的旧请求
		validIdx := len(l.requests)
		for i, reqTime := range l.requests {
			if reqTime.After(windowStart) {
				validIdx = i
				break
			}
		}
		l.requests = l.requests[validIdx:]

		// 检查是否超过阈值
		if len(l.requests) >= l.threshold {
			l.lock.Unlock()
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
		}

		// 记录当前请求
		l.requests = append(l.requests, now)
		l.lock.Unlock()

		return handler(ctx, req)
	}
}

// SlidingWindowBucketLimiter 基于分桶的滑动窗口限流器（更高效）
// 将时间窗口分成多个小桶，减少内存使用
type SlidingWindowBucketLimiter struct {
	window      time.Duration // 时间窗口大小
	bucketCount int           // 桶数量
	threshold   int           // 窗口内最大请求数
	buckets     []int         // 每个桶的计数
	bucketTime  []time.Time   // 每个桶的起始时间
	lock        sync.Mutex
}

func (l *SlidingWindowBucketLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l.lock.Lock()
		defer l.lock.Unlock()

		now := time.Now()
		bucketDuration := l.window / time.Duration(l.bucketCount)
		currentBucketIdx := int(now.UnixNano()/int64(bucketDuration)) % l.bucketCount

		// 检查当前桶是否需要重置
		if now.Sub(l.bucketTime[currentBucketIdx]) >= bucketDuration {
			l.buckets[currentBucketIdx] = 0
			l.bucketTime[currentBucketIdx] = now
		}

		// 计算滑动窗口内的总请求数
		windowStart := now.Add(-l.window)
		totalCount := 0
		for i := 0; i < l.bucketCount; i++ {
			// 只统计时间窗口内的桶
			if l.bucketTime[i].After(windowStart) || l.bucketTime[i].Equal(windowStart) {
				totalCount += l.buckets[i]
			}
		}

		// 检查是否超过阈值
		if totalCount >= l.threshold {
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
		}

		// 增加当前桶的计数
		l.buckets[currentBucketIdx]++

		return handler(ctx, req)
	}
}

//
