package tcpRouter

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net"
	"sync"
)

type TCPRouter struct {
	host     string
	ctx      map[int]*TCPDialContext
	coreFunc map[int]TCPHandlerFunc
	closer   map[int]func()

	mu sync.Mutex

	shutdown bool
	done     chan struct{}
}

func New(host string) *TCPRouter {
	return &TCPRouter{
		host:     host,
		ctx:      make(map[int]*TCPDialContext),
		coreFunc: make(map[int]TCPHandlerFunc),
		closer:   make(map[int]func()),

		mu: sync.Mutex{},

		shutdown: false,
		done:     make(chan struct{}),
	}
}

func (r *TCPRouter) TCP(port int, f TCPHandlerFunc) *TCPDialContext {
	ctx := &TCPDialContext{
		Conn:    nil,
		Context: context.Background(),
		handler: make([]TCPHandlerFunc, 0),
		index:   -1,
	}
	r.coreFunc[port] = f
	r.ctx[port] = ctx
	return ctx
}

func (r *TCPRouter) Run() {
	for k, v := range r.ctx {
		v.handler = append(v.handler, r.coreFunc[k])

		lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", r.host, k))
		if err != nil {
			hlog.Errorf("tcp router lis %s:%d error %v", r.host, k, err)
			continue
		}

		s := NewTCPServer(lis, r.ctx[k])
		r.closer[k] = s.Close
		go s.ListenAndServe()
	}

	<-r.done
}

func (r *TCPRouter) Shutdown() {
	r.mu.Lock()
	r.shutdown = true
	r.done <- struct{}{}
	r.mu.Unlock()

	fmt.Println("tcp router shutdown")
}
