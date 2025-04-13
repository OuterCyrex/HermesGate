package tcpRouter

import (
	serviceDAO "GoGateway/dao/services"
	"GoGateway/proxy"
	"GoGateway/proxy/tcp_proxy_server"
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

func (r *TCPRouter) TCP(detail *serviceDAO.ServiceDetail, f TCPHandlerFunc) *TCPDialContext {
	port := detail.Tcp.Port

	ctx := &TCPDialContext{
		Conn:    nil,
		Context: context.WithValue(context.Background(), TcpServiceDetailKey, detail),
		handler: make([]TCPHandlerFunc, 0),
		index:   -1,
	}
	r.coreFunc[port] = f
	r.ctx[port] = ctx
	return ctx
}

func (r *TCPRouter) RawTCP(port int, f TCPHandlerFunc) *TCPDialContext {
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

func (r *TCPRouter) Close(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if closer, ok := r.closer[port]; ok {
		closer()
	} else {
		hlog.Errorf("tcp router closer port %d not exist", port)
	}
}

func (r *TCPRouter) Reload(detail *serviceDAO.ServiceDetail) {
	fmt.Println("TCPReload")

	r.Close(detail.Tcp.Port)

	lb, err := proxy.ServiceBalanceHandler.GetLoadBalance(detail)
	if err != nil {
		hlog.Errorf("load balance error: %v", err)
		return
	}

	r.coreFunc[detail.Tcp.Port] = func(ctx *TCPDialContext) {
		reverseProxy := tcp_proxy_server.NewTCPReverseProxy(context.Background(), lb)
		reverseProxy.ServeTCP(ctx.Context, ctx.Conn)
	}

	r.ctx[detail.Tcp.Port].Conn = nil
	r.ctx[detail.Tcp.Port].Context = context.WithValue(context.Background(), TcpServiceDetailKey, detail)
	r.ctx[detail.Tcp.Port].index = -1

	r.ctx[detail.Tcp.Port].handler[len(r.ctx[detail.Tcp.Port].handler)-1] = r.coreFunc[detail.Tcp.Port]

	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", r.host, detail.Tcp.Port))
	if err != nil {
		hlog.Errorf("tcp router lis %s:%d error %v", r.host, detail.Tcp.Port, err)
	}

	s := NewTCPServer(lis, r.ctx[detail.Tcp.Port])
	r.closer[detail.Tcp.Port] = s.Close
	go s.ListenAndServe()
}

func (r *TCPRouter) Shutdown() {
	r.mu.Lock()
	r.shutdown = true
	r.done <- struct{}{}
	for _, closer := range r.closer {
		closer()
	}
	r.mu.Unlock()

	fmt.Println("tcp router shutdown")
}
