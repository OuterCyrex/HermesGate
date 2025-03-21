package grpc_proxy_middleware

import (
	"GoGateway/conf"
	serviceDAO "GoGateway/dao/services"
	serverProxy "GoGateway/proxy"
	"context"
	"fmt"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func NewGrpcProxyServer(detail *serviceDAO.ServiceDetail) (*grpc.Server, string) {

	addr := fmt.Sprintf("%s:%d", conf.GetConfig().ProxyServer.Host, detail.Grpc.Port)
	// 注册grpc代理服务
	s := grpc.NewServer(grpc.UnknownServiceHandler(proxy.TransparentHandler(
		func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			// 负载均衡器
			lb, err := serverProxy.ServiceBalanceHandler.GetLoadBalance(detail)
			if err != nil {
				return nil, nil, err
			}

			// 从 metadata 中获取 requestURI
			requestURI := "unknown"
			md, _ := metadata.FromIncomingContext(ctx)
			if v, ok := md[":authority"]; ok && len(v) > 0 {
				requestURI = v[0]
			}

			nextAddr, _ := lb.Get(requestURI)

			dst, err := grpc.NewClient(nextAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			return ctx, dst, err
		},
	)),
		grpc.ChainStreamInterceptor(
			GrpcLogMiddleware(),
			GrpcBlackListMiddleware(detail),
			GrpcLimitMiddleware(detail),
			GrpcFlowCountMiddleware(detail),
		))
	return s, addr
}
