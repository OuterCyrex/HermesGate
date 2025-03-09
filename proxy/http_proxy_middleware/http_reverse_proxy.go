package http_proxy_middleware

import (
	"GoGateway/conf"
	"GoGateway/pkg/status"
	"GoGateway/proxy"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/hertz-contrib/reverseproxy"
	"net/http"
	"time"
)

func HttpReverseProxyMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {

		detail := getServiceDetail(c)

		proxyAddr := fmt.Sprintf("http://%s:%d",
			conf.GetConfig().ProxyServer.Host,
			conf.GetConfig().ProxyServer.HttpPort,
		)

		proxyHandler, err := reverseproxy.NewSingleHostReverseProxy(proxyAddr,
			client.WithMaxIdleConnDuration(time.Duration(detail.LoadBalance.UpstreamIdleTimeout)*time.Second),
			client.WithDialTimeout(time.Duration(detail.LoadBalance.UpstreamConnectTimeout)*time.Second),
			client.WithClientReadTimeout(time.Duration(detail.LoadBalance.UpstreamHeaderTimeout)*time.Second),
			client.WithMaxConnsPerHost(detail.LoadBalance.UpstreamMaxIdle),
		)

		lb, err := proxy.ServiceBalanceHandler.GetLoadBalance(detail)
		if err != nil {
			status.ErrToHttpResponse(c, err)
			return
		}

		director := func(c *protocol.Request) {
			nextAddr, _ := lb.Get(string(c.RequestURI()))
			c.SetHost(nextAddr)
			c.ParseURI()
		}

		proxyHandler.SetDirector(director)
		proxyHandler.SetErrorHandler(func(c *app.RequestContext, err error) {
			c.JSON(http.StatusInternalServerError, utils.H{
				"error": err.Error(),
			})
		})

		proxyHandler.ServeHTTP(ctx, c)

		c.Next(ctx)
	}
}
