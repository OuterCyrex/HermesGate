package tcp_proxy_server

import (
	"GoGateway/proxy/load_balance"
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"net"
	"time"
)

// TCPReverseProxy 是每次TCP反向代理请求的实体/
type TCPReverseProxy struct {
	ctx             context.Context
	Addr            string
	keepAlivePeriod time.Duration
	dialTimeout     time.Duration
	DialContext     func(ctx context.Context, network, addr string) (net.Conn, error)
	OnDialError     func(conn net.Conn, err error)
}

func NewTCPReverseProxy(ctx context.Context, lb load_balance.LoadBalance, options ...ProxyOption) *TCPReverseProxy {
	nextAddr, err := lb.Get("")
	if err != nil {
		hlog.Errorf("load balance error:%v", err)
		return nil
	}
	r := &TCPReverseProxy{
		ctx:             ctx,
		Addr:            nextAddr,
		keepAlivePeriod: 2 * time.Second,
		dialTimeout:     2 * time.Second,
	}
	for _, option := range options {
		option(r)
	}

	dialer := net.Dialer{}

	if r.dialTimeout > 0 {
		dialer.Timeout = r.dialTimeout
	}

	if r.keepAlivePeriod > 0 {
		dialer.KeepAlive = r.keepAlivePeriod
	}

	r.DialContext = dialer.DialContext

	r.OnDialError = func(conn net.Conn, err error) {
		_, _ = conn.Write([]byte(err.Error()))
		hlog.Errorf("dial error: %v", err)
		_ = conn.Close()
	}

	return r
}

func (trp *TCPReverseProxy) ServeTCP(ctx context.Context, src net.Conn) {
	// 超时上下文
	var cancel context.CancelFunc
	if trp.dialTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, trp.dialTimeout)
	}

	dst, err := trp.DialContext(ctx, "tcp", trp.Addr)
	if cancel != nil {
		cancel()
	}

	if err != nil {
		trp.OnDialError(src, err)
		return
	}

	defer func() {
		_ = dst.Close()
	}()

	if trp.keepAlivePeriod > 0 {
		if tcpConn, ok := dst.(*net.TCPConn); ok {
			_ = tcpConn.SetKeepAlive(true)
			_ = tcpConn.SetKeepAlivePeriod(trp.keepAlivePeriod)
		}
	}

	errChan := make(chan error, 2)

	tcpCopy := func(src, dst net.Conn) {
		copyErr := error(nil)
		_, copyErr = io.Copy(src, dst)
		errChan <- copyErr
	}
	go tcpCopy(dst, src)
	go tcpCopy(src, dst)

	// 阻塞进程，若两个copy均结束则退出并关闭连接
	for i := 0; i < 2; i++ {
		err = <-errChan
		if err != nil && err != io.EOF {
			hlog.Errorf("tcp proxy error: %v", err)
			return
		}
	}
}
