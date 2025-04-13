package test

import (
	tcpRouter "GoGateway/proxy/tcp_router"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net"
	"testing"
)

func TestRedis(t *testing.T) {
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

func TestTCP(t *testing.T) {
	r := tcpRouter.New("127.0.0.1")
	r.RawTCP(10082, func(context *tcpRouter.TCPDialContext) {
		fmt.Println("10082")
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
