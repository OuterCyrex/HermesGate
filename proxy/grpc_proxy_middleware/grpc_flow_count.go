package grpc_proxy_middleware

import (
	serviceDAO "GoGateway/dao/services"
	redisCounter "GoGateway/proxy/redis_counter"
	"google.golang.org/grpc"
)

func GrpcFlowCountMiddleware(detail *serviceDAO.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		counter := redisCounter.ServiceFlowCountHandler.GetCounter(detail.Info.ServiceName)

		counter.Increase()

		return handler(srv, ss)
	}
}
