package http_proxy_middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"strings"
)

func HttpRewriteHeaderMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		detail := getServiceDetail(c)

		if detail.Http.HeaderTransfer == "" {
			c.Next(ctx)
			return
		}

		transfers := strings.Split(detail.Http.HeaderTransfer, ",")

		SetMap := make(map[string]string)
		var DelList []string

		for _, t := range transfers {
			ops := strings.Split(t, " ")
			if !checkTransferValid(ops, 2, c) {
				return
			}
			switch ops[0] {
			case "add":
				if !checkTransferValid(ops, 3, c) {
					return
				}
				SetMap[ops[1]] = ops[2]
			case "edit":
				if !checkTransferValid(ops, 3, c) {
					return
				}
				SetMap[ops[1]] = ops[2]
			case "del":
				if !checkTransferValid(ops, 2, c) {
					return
				}
				DelList = append(DelList, ops[1])
			default:
				continue
			}
		}

		c.Request.SetHeaders(SetMap)

		for _, d := range DelList {
			c.Request.Header.Del(d)
		}

		c.Next(ctx)
	}
}

func checkTransferValid(transfers []string, limit int, c *app.RequestContext) bool {
	if len(transfers) < limit {
		c.JSON(http.StatusInternalServerError, utils.H{
			"message": "Invalid header transfer format",
		})
		c.Abort()
		return false
	}
	return true
}
