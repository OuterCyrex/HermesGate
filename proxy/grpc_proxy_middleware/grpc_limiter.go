package grpc_proxy_middleware

import (
	serviceDAO "GoGateway/dao/services"
	"GoGateway/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"strings"
)

func GrpcLimitMiddleware(detail *serviceDAO.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if detail.AccessControl.ServiceFlowLimit != 0 {
			limiter := proxy.ServiceLimitHandler.GetServerLimiter(detail.Info.ServiceName, float64(detail.AccessControl.ServiceFlowLimit))
			if !limiter.Allow() {
				return status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
			}
		}

		clientIP := getIPInfo(ss)

		if detail.AccessControl.ClientIPFlowLimit != 0 {
			limiter := proxy.ServiceLimitHandler.GetClientLimiter(detail.Info.ServiceName, clientIP, float64(detail.AccessControl.ClientIPFlowLimit))
			if !limiter.Allow() {
				return status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
			}
		}

		return handler(srv, ss)
	}
}

func getIPInfo(ss grpc.ServerStream) string {
	clientIP := "unknown"
	if p, ok := peer.FromContext(ss.Context()); ok {
		addr := p.Addr.String()
		if v := strings.Split(addr, ":"); len(v) > 1 {
			clientIP = v[0]
		}
	}
	return clientIP
}
