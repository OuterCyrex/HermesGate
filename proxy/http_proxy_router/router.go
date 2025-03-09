package http_proxy_router

import (
	"GoGateway/biz/middleware"
	"GoGateway/conf"
	"GoGateway/pkg"
	"GoGateway/proxy/http_proxy_middleware"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func InitHttpProxyRouter() func(ctx context.Context) (err error) {
	addr := fmt.Sprintf("%s:%d", conf.GetConfig().ProxyServer.Host, conf.GetConfig().ProxyServer.HttpPort)
	s := server.Default(server.WithHostPorts(addr))

	s.Use(pkg.GetCors())
	s.Use(middleware.AccessLog())
	s.Use(
		http_proxy_middleware.HttpAccessMiddleware(),
		http_proxy_middleware.HttpReverseProxyMiddleware(),
	)

	s.Spin()

	return s.Shutdown
}
