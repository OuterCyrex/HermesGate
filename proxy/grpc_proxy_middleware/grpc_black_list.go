package grpc_proxy_middleware

import (
	serviceDAO "GoGateway/dao/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func GrpcBlackListMiddleware(detail *serviceDAO.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var blackList []string
		var whiteList []string
		if detail.AccessControl.BlackList != "" {
			blackList = strings.Split(detail.AccessControl.BlackList, ",")
		}

		if detail.AccessControl.WhiteList != "" {
			whiteList = strings.Split(detail.AccessControl.WhiteList, ",")
		}

		ip := getIPInfo(ss)

		// 白名单优先
		if detail.AccessControl.OpenAuth == 1 && len(detail.AccessControl.WhiteList) > 0 {
			for _, w := range whiteList {
				if w == ip {
					return handler(srv, ss)
				}
			}
			return status.Errorf(codes.PermissionDenied, "access control white list not matches")
		}

		// 若白名单为空则使用黑名单
		if detail.AccessControl.OpenAuth == 1 && len(detail.AccessControl.WhiteList) == 0 && len(detail.AccessControl.BlackList) > 0 {
			for _, w := range blackList {
				if w == ip {
					return status.Errorf(codes.PermissionDenied, "access control black list matches")
				}
			}
			return handler(srv, ss)
		}

		// 若未开启权限验证或黑白名单均为空则放行
		return handler(srv, ss)
	}
}
