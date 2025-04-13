package tcp_proxy_middleware

import (
	"GoGateway/proxy"
	tcpRouter "GoGateway/proxy/tcp_router"
)

func TcpLimitMiddleware() tcpRouter.TCPHandlerFunc {
	return func(c *tcpRouter.TCPDialContext) {
		detail := c.GetDetail()

		if detail.AccessControl.ServiceFlowLimit != 0 {
			limiter := proxy.ServiceLimitHandler.GetServerLimiter(detail.Info.ServiceName, float64(detail.AccessControl.ServiceFlowLimit))
			if !limiter.Allow() {
				c.Abort()
				return
			}
		}

		if detail.AccessControl.ClientIPFlowLimit != 0 {
			limiter := proxy.ServiceLimitHandler.GetClientLimiter(detail.Info.ServiceName, c.ClientIP(), float64(detail.AccessControl.ClientIPFlowLimit))
			if !limiter.Allow() {
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
