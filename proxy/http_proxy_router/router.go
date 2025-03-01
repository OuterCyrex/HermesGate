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
	addr := fmt.Sprintf("%s:%d", conf.GetConfig().ProxyServer.HttpProxy.Host, conf.GetConfig().ProxyServer.HttpProxy.Port)
	s := server.Default(server.WithHostPorts(addr))

	s.Use(pkg.GetCors())
	s.Use(middleware.AccessLog())
	s.Use(
		http_proxy_middleware.HttpAccessMiddleware(),
	)

	s.Spin()

	return s.Shutdown
}
