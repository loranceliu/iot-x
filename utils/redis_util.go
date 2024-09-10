package utils

import (
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type RedisUtil struct {
	client *redis.Client
}

func NewRedisUtil(client *redis.Client) *RedisUtil {
	return &RedisUtil{client: client}
}

// Set 将指定值设置到redis中，使用json的序列化方式
func (u *RedisUtil) SetObj(key string, value interface{}, duration time.Duration) error {
	bytes, err := jsoniter.Marshal(value)
	if err != nil {
		return err
	}

	err = u.client.Set(key, bytes, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get 从redis中读取指定值，使用json的反序列化方式
func (u *RedisUtil) GetObj(key string, value interface{}) error {
	bytes, err := u.client.Get(key).Bytes()
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(bytes, value)
	if err != nil {
		return err
	}
	return nil
}

func (u *RedisUtil) Del(key string) error {
	_, err := u.client.Del(key).Result()
	return err
}
