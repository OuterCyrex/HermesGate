package tcp_proxy_middleware

import (
	redisCounter "GoGateway/proxy/redis_counter"
	tcpRouter "GoGateway/proxy/tcp_router"
)

func TcpFlowCounterMiddleware() tcpRouter.TCPHandlerFunc {
	return func(c *tcpRouter.TCPDialContext) {

		detail := c.GetDetail()

		counter := redisCounter.ServiceFlowCountHandler.GetCounter(detail.Info.ServiceName)

		counter.Increase()

		c.Next()
	}
}
