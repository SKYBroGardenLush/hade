package contract

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	*redis.Options
}

func (redis *RedisConfig) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v_%v", redis.Addr, redis.DB, redis.Username, redis.Password)
}
