package storage

import (
	"time"
)

type RedisCacheMock struct {
	m map[string]interface{}
}

func (r *RedisCacheMock) Set(key string, value interface{}, expiration time.Duration) error {
	r.m[key] = value
	return nil
}

func BuildTestRedisClient() {
	mock := &RedisCacheMock{make(map[string]interface{})}
	redisClient = mock
}
