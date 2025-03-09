package test

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"testing"
)

func TestReverseProxy(t *testing.T) {
	h := server.New(server.WithHostPorts("127.0.0.1:8080"))

	h.GET("/user", func(ctx context.Context, c *app.RequestContext) {

		c.JSON(http.StatusOK, utils.H{
			"host":       string(c.Request.Host()),
			"RequestURI": string(c.Request.RequestURI()),
			"Path":       string(c.Request.Path()),
			"Schema":     string(c.Request.Scheme()),
			"uri":        c.Request.URI().String(),
		})
	})

	h.Spin()
}
