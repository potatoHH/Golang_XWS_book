package cache

import (
	"Book_Exp/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// var ErrKeyNotExist = errors.New("key 不存在")
var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
	Delete(ctx context.Context, id int64) error
}

type RedisUserCache struct {
	client     redis.Cmdable // 命令  可以传单机的redis or 传cluster的redis
	expiration time.Duration // 过期时间
}

// 依赖注入
// A用到了B,B一定是接口
// A用到了B,B一定A的字段
// A用到了B,A绝对不初始化B,而是外面注入
func NewUserCache(client redis.Cmdable) UserCache { // 这里返回的是一个指针
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}
func (cache *RedisUserCache) Delete(ctx context.Context, id int64) error {
	return cache.client.Del(ctx, cache.Key(id)).Err()
}

// 只有error为nil,就认为缓存里有数据
// 如果没有数据,返回一个特定的error
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.Key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err

}
func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.Key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}
func (cache *RedisUserCache) Key(id int64) string {
	return fmt.Sprintf("user:info %d", id)
}
