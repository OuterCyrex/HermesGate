package tcpRouter

import (
	"context"
	"math"
	"net"
)

const abortIndex int8 = math.MaxInt8 / 2

type TCPHandlerFunc func(*TCPDialContext)
type TCPDialContext struct {
	conn    net.Conn
	context context.Context
	handler []TCPHandlerFunc
	index   int8
}

func (c *TCPDialContext) Get(key string) interface{} {
	return c.context.Value(key)
}

func (c *TCPDialContext) Set(key string, value interface{}) {
	c.context = context.WithValue(c.context, key, value)
}

func (c *TCPDialContext) ClientIP() string {
	return c.conn.RemoteAddr().String()
}

func (c *TCPDialContext) Use(handlers ...TCPHandlerFunc) {
	c.Reset()
	c.handler = append(c.handler, handlers...)
}

func (c *TCPDialContext) Next() {
	for c.index < int8(len(c.handler))-1 && c.index < abortIndex {
		c.index++
		c.handler[c.index](c)
	}
}

func (c *TCPDialContext) Abort() {
	c.index = abortIndex
}

func (c *TCPDialContext) Reset() {
	c.index = -1
}

func (c *TCPDialContext) IsAborted() bool {
	return c.index >= abortIndex
}
