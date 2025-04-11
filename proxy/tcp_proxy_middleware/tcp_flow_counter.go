package tcp_proxy_middleware

import (
	serviceDAO "GoGateway/dao/services"
	redisCounter "GoGateway/proxy/redis_counter"
	tcpRouter "GoGateway/proxy/tcp_router"
)

func TcpFlowCounterMiddleware(detail *serviceDAO.ServiceDetail) tcpRouter.TCPHandlerFunc {
	return func(c *tcpRouter.TCPDialContext) {

		counter := redisCounter.ServiceFlowCountHandler.GetCounter(detail.Info.ServiceName)

		counter.Increase()

		c.Next()
	}
}
