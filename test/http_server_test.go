package test

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"testing"
)

func TestReverseProxy1(t *testing.T) {
	h := server.New(server.WithHostPorts("127.0.0.1:8081"))

	h.GET("/user", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, utils.H{
			"msg":     "8081收到",
			"headers": c.Request.Header.Get("Authorization"),
		})
	})

	h.Spin()
}

func TestReverseProxy2(t *testing.T) {
	h := server.New(server.WithHostPorts("127.0.0.1:8082"))

	h.GET("/user", func(ctx context.Context, c *app.RequestContext) {

		c.JSON(http.StatusOK, utils.H{
			"msg":     "8082收到",
			"headers": c.Request.Header.Get("Authorization"),
		})
	})

	h.Spin()
}
