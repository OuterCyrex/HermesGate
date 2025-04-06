package tcpRouter

import (
	"GoGateway/conf"
	serviceConsts "GoGateway/pkg/consts/service"
	"GoGateway/proxy"
	"GoGateway/proxy/tcp_proxy_server"
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var TcpRouter *TCPRouter

func InitTCPRouter(manager *proxy.ServiceManager) {
	TcpRouter = New(conf.GetConfig().ProxyServer.Host)

	for _, v := range manager.ServiceMap {
		if v.Info.LoadType == serviceConsts.ServiceLoadTypeTCP {
			lb, err := proxy.ServiceBalanceHandler.GetLoadBalance(v)
			if err != nil {
				hlog.Errorf("load balance error: %v", err)
				return
			}
			reverseProxy := tcp_proxy_server.NewTCPReverseProxy(context.Background(), lb)
			TcpRouter.TCP(v.Tcp.Port, func(ctx *TCPDialContext) {
				reverseProxy.ServeTCP(ctx.context, ctx.conn)
			})
		}
	}

	go TcpRouter.Run()
}
