package grpc_proxy_server

import (
	serviceDAO "GoGateway/dao/services"
	serviceConsts "GoGateway/pkg/consts/service"
	"GoGateway/proxy"
	"GoGateway/proxy/grpc_proxy_middleware"
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type GrpcServerManager struct {
	serverMap map[string]*grpcServer
	manager   *proxy.ServiceManager
	lock      sync.Mutex
}

type grpcServer struct {
	server *grpc.Server
	addr   string
}

func NewGrpcServerManager(manager *proxy.ServiceManager) *GrpcServerManager {
	if manager == nil {
		return nil
	}

	m := &GrpcServerManager{
		serverMap: make(map[string]*grpcServer),
		manager:   manager,
		lock:      sync.Mutex{},
	}

	for _, detail := range manager.ServiceMap {
		if detail.Info.LoadType == serviceConsts.ServiceLoadTypeGRPC {
			server, addr := grpc_proxy_middleware.NewGrpcProxyServer(detail)
			m.lock.Lock()
			m.serverMap[detail.Info.ServiceName] = &grpcServer{
				server: server,
				addr:   addr,
			}
			m.lock.Unlock()
		}
	}

	return m
}

func (s *GrpcServerManager) Serve() error {
	for _, i := range s.serverMap {
		go func() {
			hlog.Infof("grpc proxy server serve on %s",
				i.addr,
			)

			lis, lisErr := net.Listen("tcp", i.addr)

			serverErr := i.server.Serve(lis)

			if err := errors.Join(lisErr, serverErr); err != nil {
				panic(err)
			}
		}()
	}
	return nil
}

func (s *GrpcServerManager) Reload(detail *serviceDAO.ServiceDetail) {
	s.lock.Lock()
	// 关闭先前的grpc服务
	if v, ok := s.serverMap[detail.Info.ServiceName]; ok {
		v.server.GracefulStop()
	}

	server, addr := grpc_proxy_middleware.NewGrpcProxyServer(detail)
	s.serverMap[detail.Info.ServiceName] = &grpcServer{
		server: server,
		addr:   addr,
	}

	s.lock.Unlock()

	go func() {
		lis, lisErr := net.Listen("tcp", addr)

		serverErr := server.Serve(lis)

		if err := errors.Join(lisErr, serverErr); err != nil {
			hlog.Errorf("grpc proxy server serve error: %v", err)
		}
	}()
}
