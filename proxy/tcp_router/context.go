package tcpRouter

import (
	"context"
	"math"
	"net"
	"strings"
)

const abortIndex int8 = math.MaxInt8 / 2

type TCPHandlerFunc func(*TCPDialContext)
type TCPDialContext struct {
	Conn    net.Conn
	Context context.Context
	handler []TCPHandlerFunc
	index   int8
}

func (c *TCPDialContext) Get(key string) interface{} {
	return c.Context.Value(key)
}

func (c *TCPDialContext) Set(key string, value interface{}) {
	c.Context = context.WithValue(c.Context, key, value)
}

func (c *TCPDialContext) ClientIP() string {
	addr := c.Conn.RemoteAddr().String()
	addrSlice := strings.Split(addr, ":")
	if len(addrSlice) == 2 {
		return addrSlice[0]
	} else {
		return ""
	}
}

func (c *TCPDialContext) Use(handlers ...TCPHandlerFunc) {
	c.Reset()
	c.handler = append(c.handler, handlers...)
}

func (c *TCPDialContext) Write(b []byte) (n int, err error) {
	return c.Conn.Write(b)
}

func (c *TCPDialContext) Read(b []byte) (n int, err error) {
	return c.Conn.Read(b)
}

func (c *TCPDialContext) Next() {
	c.index++
	for c.index < int8(len(c.handler)) && c.index < abortIndex {
		c.handler[c.index](c)
		c.index++
	}
}

func (c *TCPDialContext) Abort() {
	_ = c.Conn.Close()
	c.index = abortIndex
}

func (c *TCPDialContext) Reset() {
	c.index = -1
}

func (c *TCPDialContext) IsAborted() bool {
	return c.index >= abortIndex
}
