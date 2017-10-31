package storage

import (
	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
)

func GetRedisClient() (*redis.Client, error) {
	if redisClient != nil {
		return redisClient, nil
	}
	//TODO(pavel): move to the configiration
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}
