package test

import (
	tcpRouter "GoGateway/proxy/tcp_router"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

func TestTCP(t *testing.T) {
	r := tcpRouter.New("127.0.0.1")

	var count atomic.Int64

	count.Store(0)

	r.TCP(9910, func(context *tcpRouter.TCPDialContext) {
		count.Add(1)
		log.Println(count.Load())
	})

	r.Run()

	r.Shutdown()
}

func TestTCPRouter(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			log.Println("send")

			conn, err := net.Dial("tcp", "127.0.0.1:8971")
			if err != nil {
				fmt.Println(err)
			}

			conn.Write([]byte("hello"))
			conn.Close()
		}()
	}

	time.Sleep(time.Second * 3)
}
