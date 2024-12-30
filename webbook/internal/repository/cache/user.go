package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"jikeshijian_go/webbook/internal/domain"
	"time"
)

// redis.Nil 就是 key 不存在
var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

// NewRedisUserCache 构造函数
// A 用到了B ，B 一定是接口 ==》 面向接口编程
// A 用到了B ，B 一定是A的字段 ==》 规避包变量，包方法，都非常缺乏扩展性
// A 用到了B ，A绝不初始化B，而是外面注入 ==》 保持依赖注入（DI和IOC）
func NewRedisUserCache(cmd redis.Cmdable, expiration time.Duration) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: expiration,
	}
}

func (c *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.key(uid)
	// 我假定这个地方用 JSON 存储
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	//if err != nil {
	//	return domain.User{}, err
	//}
	//return u, nil
	return u, err
}

func (c *RedisUserCache) Set(ctx context.Context, du domain.User) error {
	key := c.key(du.Id)
	// 我假定这个地方用 JSON
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisUserCache) key(uid int64) string {
	// user-info-
	// user.info.
	// user/info/
	// user_info_
	return fmt.Sprintf("user:info:%d", uid)
}
