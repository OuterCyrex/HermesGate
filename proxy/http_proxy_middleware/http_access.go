package http_proxy_middleware

import (
	serviceDAO "GoGateway/dao/services"
	"GoGateway/pkg/status"
	"GoGateway/proxy"
	"context"
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
		setServiceDetail(c, detail)
		c.Next(ctx)
	}
}

func setServiceDetail(c *app.RequestContext, detail *serviceDAO.ServiceDetail) {
	c.Set("service", detail)
}

func getServiceDetail(c *app.RequestContext) *serviceDAO.ServiceDetail {
	if detail, ok := c.MustGet("service").(*serviceDAO.ServiceDetail); ok {
		return detail
	} else {
		return nil
	}
}
