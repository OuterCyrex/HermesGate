package middleware

import (
	"GoGateway/pkg"
	redisConsts "GoGateway/pkg/consts/redis"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strconv"
)

func HotReloadMiddleware(ctx context.Context, app *app.RequestContext) {
	app.Next(ctx)

	id, err := strconv.Atoi(app.Param("id"))
	if err != nil {
		hlog.Errorf("hot reload failed: %v", err.Error())
	}

	err = PubSend(id)

	if err != nil {
		hlog.Errorf("hot reload failed: %v", err.Error())
	}
}

func PubSend(id int) error {
	conn := pkg.GetRedis()
	defer func() {
		_ = conn.Close()
	}()

	_, err := conn.Do("PUBLISH", redisConsts.RedisChannelKey, id)
	if err != nil {
		hlog.Errorf("hot reload failed: %v", err.Error())
		return err
	}
	return nil
}
