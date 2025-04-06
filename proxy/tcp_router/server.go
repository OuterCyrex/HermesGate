package tcpRouter

import (
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net"
)

var (
	AlreadyCloseErr = errors.New("server already close")
)

type TCPServer struct {
	lis      net.Listener
	ctx      TCPDialContext
	done     chan struct{}
	shutdown bool
}

func (ts *TCPServer) Close() {
	if ts.shutdown {
		return
	}
	ts.shutdown = true
	ts.done <- struct{}{}
	err := ts.lis.Close()
	if err != nil {
		hlog.Errorf("tcp server close error %v", err)
	}
	close(ts.done)
}

func NewTCPServer(lis net.Listener, ctx *TCPDialContext) *TCPServer {
	return &TCPServer{
		lis:      lis,
		ctx:      *ctx,
		done:     make(chan struct{}),
		shutdown: false,
	}
}

func (ts *TCPServer) ListenAndServe() error {
	if ts.shutdown {
		return AlreadyCloseErr
	}

	go func() {
		for {
			select {
			case <-ts.done:
				hlog.Debug("tcp server exit")
				return
			default:
			}

			conn, err := ts.lis.Accept()
			if err != nil {
				hlog.Errorf("tcp server accept error %v", err)
			}

			go func() {
				ts.ctx.conn = conn
				ts.ctx.Reset()
				ts.ctx.Next()
			}()
		}
	}()

	return nil
}
