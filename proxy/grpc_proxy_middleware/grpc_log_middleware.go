package grpc_proxy_middleware

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

func GrpcLogMiddleware() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		clientIP := "unknown"
		if p, ok := peer.FromContext(ctx); ok {
			clientIP = p.Addr.String()
		}

		host := "unknown"
		md, _ := metadata.FromIncomingContext(ctx)
		if v, ok := md[":authority"]; ok && len(v) > 0 {
			host = v[0]
		}

		startTime := time.Now()

		err := handler(srv, ss)

		endTime := time.Now()
		duration := endTime.Sub(startTime)
		hlog.Debugf("status=%d full_method=%s client_ip=%s host=%s cost=%s",
			getStatus(err),
			info.FullMethod,
			clientIP,
			host,
			timeFormat(duration),
		)
		return err
	}
}

func timeFormat(d time.Duration) string {
	if d.Microseconds() < 1000 {
		return strconv.Itoa(int(d.Microseconds())) + "μs"
	} else {
		micro := float64(d.Microseconds())
		return fmt.Sprintf("%.2f", micro/1000) + "ms"
	}
}

func getStatus(err error) int {
	if err == nil {
		return 200
	}
	switch status.Code(err) {
	case codes.OK:
		return 200
	case codes.NotFound:
		return 404
	case codes.Unimplemented:
		return 501
	case codes.Unavailable:
		return 503
	case codes.Internal:
		return 500
	default:
		return 500
	}
}
