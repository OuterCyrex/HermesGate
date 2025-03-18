package pkg

import (
	"GoGateway/conf"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

func InitRedis() {
	redisPool = &redis.Pool{
		MaxIdle:     180,
		MaxActive:   200,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d",
				conf.GetConfig().DashBoard.Redis.Host,
				conf.GetConfig().DashBoard.Redis.Port,
			),
				redis.DialPassword(conf.GetConfig().DashBoard.Redis.Password),
				redis.DialDatabase(conf.GetConfig().DashBoard.Redis.Database),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	conn := redisPool.Get()
	defer func() {
		_ = conn.Close()
	}()
	_, err := conn.Do("ping")
	if err != nil {
		panic(err)
	}
}

var redisPool *redis.Pool

func GetRedis() redis.Conn {
	return redisPool.Get()
}
