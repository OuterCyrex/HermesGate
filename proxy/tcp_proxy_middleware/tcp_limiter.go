package tcp_proxy_middleware

import (
	serviceDAO "GoGateway/dao/services"
	"GoGateway/proxy"
	tcpRouter "GoGateway/proxy/tcp_router"
)

func TcpLimitMiddleware(detail *serviceDAO.ServiceDetail) tcpRouter.TCPHandlerFunc {
	return func(c *tcpRouter.TCPDialContext) {
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
