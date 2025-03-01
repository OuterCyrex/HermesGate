package http_proxy_middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func HttpReverseProxyMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {

		c.Next(ctx)
	}
}
