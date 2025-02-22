// Code generated by hertz generator.

package admin

import (
	"GoGateway/biz/middleware"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/cookie"
	"github.com/hertz-contrib/sessions/redis"
	"sync"
)

var (
	once  sync.Once
	store redis.Store
)

func getStore() redis.Store {
	once.Do(func() {
		store = cookie.NewStore([]byte("secret"))
		// In Redis way:
		// store, _ = redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	})
	return store
}

func rootMw() []app.HandlerFunc {
	return nil
}

func _adminMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{sessions.New("admin", getStore())}
}

func _adminloginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _admininfoMw() []app.HandlerFunc {
	return []app.HandlerFunc{middleware.SessionAuthMiddleware}
}

func _adminlogoutMw() []app.HandlerFunc {
	return []app.HandlerFunc{middleware.SessionAuthMiddleware}
}

func _changepasswordMw() []app.HandlerFunc {
	return []app.HandlerFunc{middleware.SessionAuthMiddleware}
}
