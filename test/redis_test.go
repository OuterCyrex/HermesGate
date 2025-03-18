package test

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	redisPool := &redis.Pool{
		MaxIdle:     180,
		MaxActive:   200,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d",
				"127.0.0.1",
				6379,
			),
				redis.DialPassword(""),
				redis.DialDatabase(0),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	conn := redisPool.Get()
	fmt.Println(conn.Do("GET", "redis_counter_HttpTestServer_20250318"))
}
