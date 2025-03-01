package http_proxy_middleware

import (
	"GoGateway/pkg/status"
	"GoGateway/proxy"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
)

func HttpAccessMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		detail, err := proxy.ServiceManagerHandler.GetHttpDetail(c)
		if err != nil {
			status.ErrToHttpResponse(c, err)
			c.Abort()
			return
		}
		fmt.Println(detail)
		c.Next(ctx)
	}
}
