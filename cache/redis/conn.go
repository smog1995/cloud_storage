package redis

import (
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"time"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
	redisPass = "testupload"
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				zap.S().Error(err.Error())
				return nil, err
			}

			//访问认证
			if _, err = c.Do("AUTH", redisPass); err != nil {
				zap.S().Error(err.Error())
				c.Close()
				return nil, err
			}
			return c, nil
		},
		// 每分钟检查redis连接的可用性
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute { //小于一分钟返回nil
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()

	_, err := pool.Get().Do("KEYS", "*")
	if err != nil {
		zap.S().Error(err.Error())
	}
}

func RedisPool() *redis.Pool {
	return pool
}
