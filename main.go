// Code generated by hertz generator.

package main

import (
	"GoGateway/biz/middleware"
	"GoGateway/conf"
	"GoGateway/dao"
	"GoGateway/pkg"
	"GoGateway/proxy"
	"GoGateway/proxy/http_proxy_router"
	"context"
	"flag"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"os"
	"os/signal"
	"syscall"
)

var (
	endpoint = flag.String("endpoint", "proxy-server", "input endpoint like dashboard or proxy-server")
)

func main() {
	// parse param
	flag.Parse()

	switch *endpoint {
	case "dashboard":
		dashBoardEndPoint()
	case "proxy-server":
		proxyServerEndPoint()
	default:
		flag.Usage()
	}
}

func dashBoardEndPoint() {
	hlog.SetLevel(hlog.LevelDebug)
	hlog.SetOutput(os.Stdout)

	dao.InitDB(dao.DefaultDSN())

	addr := fmt.Sprintf("%s:%d", conf.GetConfig().DashBoard.Host, conf.GetConfig().DashBoard.Port)

	h := server.Default(server.WithHostPorts(addr))

	h.Use(pkg.GetCors())
	h.Use(middleware.AccessLog())

	register(h)
	h.Spin()
}

func proxyServerEndPoint() {
	hlog.SetLevel(hlog.LevelDebug)
	hlog.SetOutput(os.Stdout)

	dao.InitDB(dao.DefaultDSN())

	err := proxy.ServiceManagerHandler.LoadOnce()
	if err != nil {
		hlog.Fatalf("load service manager error %v", err)
	}

	type Closer func(ctx context.Context) (err error)

	var closers []Closer

	go func() {
		closers = append(closers, http_proxy_router.InitHttpProxyRouter())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	hlog.Debug("shutdown proxy servers received\n")

	ctx := context.Background()

	for _, c := range closers {
		err := c(ctx)
		if err != nil {
			hlog.Errorf("shutdown proxy servers failed: %v", err.Error())
		}
	}

	hlog.Debug("shutdown proxy servers success\n")
}
