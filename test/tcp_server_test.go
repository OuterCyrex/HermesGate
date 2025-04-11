package test

import (
	tcpRouter "GoGateway/proxy/tcp_router"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net"
	"testing"
)

func TestRedis(t *testing.T) {
	for {
		c, err := redis.Dial("tcp", "127.0.0.1:8971")
		if err != nil {
			t.Error(err)
		}
		r, err := c.Do("Ping")
		if err != nil {
			t.Error(err)
		}

		fmt.Println(r.(string))
	}
}

func TestTCP(t *testing.T) {
	r := tcpRouter.New("127.0.0.1")
	r.TCP(10082, func(context *tcpRouter.TCPDialContext) {
		fmt.Println("core")
	}).Use(func(c *tcpRouter.TCPDialContext) {
		fmt.Println("111")
		c.Next()
	}, func(c *tcpRouter.TCPDialContext) {
		fmt.Println("222")
		c.Abort()
	}, func(c *tcpRouter.TCPDialContext) {
		c.Next()
		fmt.Println("333")
	})

	r.Run()
}

func TestClient(t *testing.T) {
	c, err := net.Dial("tcp", "127.0.0.1:10082")
	if err != nil {
		t.Error(err)
	}

	c.Write([]byte("hello world"))
}
