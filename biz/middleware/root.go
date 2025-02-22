package middleware

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strconv"
	"time"
)

func AccessLog() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		end := time.Now()
		latency := end.Sub(start)
		hlog.CtxDebugf(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s cost=%s",
			c.Response.StatusCode(),
			c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host(), timeFormat(latency))
	}
}

func timeFormat(d time.Duration) string {
	if d.Microseconds() < 1000 {
		return strconv.Itoa(int(d.Microseconds())) + "μs"
	} else {
		micro := float64(d.Microseconds())
		return fmt.Sprintf("%.2f", micro/1000) + "ms"
	}
}
