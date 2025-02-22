package middleware

import (
	adminDAO "GoGateway/dao/admin"
	"GoGateway/pkg/consts/session"
	"GoGateway/pkg/status"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
	"net/http"
)

func SessionAuthMiddleware(ctx context.Context, app *app.RequestContext) {
	ses := sessions.Default(app)
	sesInfo := ses.Get(sessionKey.AdminSessionInfoKey)

	v, ok := sesInfo.([]byte)

	if !ok {
		app.JSON(http.StatusForbidden, status.NewErrorResponse("请先登录"))
		app.Abort()
		return
	}

	userInfo := adminDAO.AdminSessionInfo{}

	err := userInfo.Unmarshal(v)

	if err != nil {
		hlog.Errorf("unmarshal adminInfo failed: %v", err.Error())
		app.JSON(http.StatusInternalServerError, status.NewErrorResponse("服务器内部错误"))
		app.Abort()
		return
	}

	app.Set(sessionKey.AdminSessionInfoKey, userInfo)

	app.Next(ctx)
}
