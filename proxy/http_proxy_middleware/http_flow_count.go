package http_proxy_middleware

import (
	"GoGateway/proxy/redis_counter"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func HttpFlowCountMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		detail := getServiceDetail(c)

		counter := redisCounter.ServiceFlowCountHandler.GetCounter(detail.Info.ServiceName)

		counter.Increase()

		c.Next(ctx)
	}
}
