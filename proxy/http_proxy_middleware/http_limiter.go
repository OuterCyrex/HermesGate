package http_proxy_middleware

import (
	"GoGateway/proxy"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
)

func HttpFlowLimiterMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		detail := getServiceDetail(c)

		if detail.AccessControl.ServiceFlowLimit != 0 {
			limiter := proxy.ServiceLimitHandler.GetServerLimiter(detail.Info.ServiceName, float64(detail.AccessControl.ServiceFlowLimit))
			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, utils.H{
					"message": "rate limit exceeded",
				})
				c.Abort()
			}
		}

		if detail.AccessControl.ClientIPFlowLimit != 0 {
			limiter := proxy.ServiceLimitHandler.GetClientLimiter(detail.Info.ServiceName, c.ClientIP(), float64(detail.AccessControl.ClientIPFlowLimit))
			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, utils.H{
					"message": "rate limit exceeded",
				})
				c.Abort()
			}
		}

		c.Next(ctx)
	}
}
