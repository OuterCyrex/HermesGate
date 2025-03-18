package http_proxy_middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"strings"
)

func HttpBlackListMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		detail := getServiceDetail(c)

		var blackList []string
		var whiteList []string
		if detail.AccessControl.BlackList != "" {
			blackList = strings.Split(detail.AccessControl.BlackList, ",")
		}

		if detail.AccessControl.WhiteList != "" {
			whiteList = strings.Split(detail.AccessControl.WhiteList, ",")
		}

		// 白名单优先
		if detail.AccessControl.OpenAuth == 1 && len(detail.AccessControl.WhiteList) > 0 {
			for _, w := range whiteList {
				if w == c.ClientIP() {
					c.Next(ctx)
					return
				}
			}
			c.JSON(http.StatusForbidden, utils.H{
				"forbidden": "access control white list not matches",
			})
			c.Abort()
			return
		}

		// 若白名单为空则使用黑名单
		if detail.AccessControl.OpenAuth == 1 && len(detail.AccessControl.WhiteList) == 0 && len(detail.AccessControl.BlackList) > 0 {
			for _, w := range blackList {
				if w == c.ClientIP() {
					c.JSON(http.StatusForbidden, utils.H{
						"forbidden": "access control black list matches",
					})
					c.Abort()
					return
				}
			}
			c.Next(ctx)
			return
		}

		// 若未开启权限验证或黑白名单均为空则放行
		c.Next(ctx)
	}
}
