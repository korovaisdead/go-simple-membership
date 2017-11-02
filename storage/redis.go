package storage

import (
	"github.com/go-redis/redis"
	"github.com/korovaisdead/go-simple-membership/config"
	"time"
)

var (
	redisClient RedisCache
)

type RedisCache interface {
	Set(key string, value interface{}, expiration time.Duration) error
}

type RedisCacheImpl struct {
	redisClient *redis.Client
}

func (r *RedisCacheImpl) Set(key string, value interface{}, expiration time.Duration) error {
	_, err := r.redisClient.Set(key, value, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetRedisClient() RedisCache {
	if redisClient == nil {
		panic("You should build the redis client before using it!")
	}

	return redisClient
}

func BuildRedisClient() error {
	c, _ := config.GetConfig()
	if redisClient != nil {
		return nil
	}

	r := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host + c.Redis.Port,
		Password: "",
		DB:       c.Redis.Database,
	})

	_, err := r.Ping().Result()
	if err != nil {
		return err
	}
	redisClient = &RedisCacheImpl{r}
	return nil
}
