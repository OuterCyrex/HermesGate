package tcp_proxy_middleware

import (
	"GoGateway/conf"
	serviceConsts "GoGateway/pkg/consts/service"
	"GoGateway/proxy"
	"GoGateway/proxy/tcp_proxy_server"
	"GoGateway/proxy/tcp_router"
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var TcpRouter *tcpRouter.TCPRouter

func InitTCPRouter(manager *proxy.ServiceManager) {
	TcpRouter = tcpRouter.New(conf.GetConfig().ProxyServer.Host)

	for _, v := range manager.ServiceMap {
		if v.Info.LoadType == serviceConsts.ServiceLoadTypeTCP {
			lb, err := proxy.ServiceBalanceHandler.GetLoadBalance(v)
			if err != nil {
				hlog.Errorf("load balance error: %v", err)
				return
			}
			TcpRouter.TCP(v.Tcp.Port, func(ctx *tcpRouter.TCPDialContext) {
				reverseProxy := tcp_proxy_server.NewTCPReverseProxy(context.Background(), lb)
				reverseProxy.ServeTCP(ctx.Context, ctx.Conn)
			}).Use(
				TcpBlackListMiddleware(v),
				TcpLimitMiddleware(v),
				TcpFlowCounterMiddleware(v),
			)
		}
	}

	go TcpRouter.Run()
}
