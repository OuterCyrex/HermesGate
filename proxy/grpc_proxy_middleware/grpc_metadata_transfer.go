package grpc_proxy_middleware

import (
	serviceDAO "GoGateway/dao/services"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func GrpcRewriteMetadataMiddleware(detail *serviceDAO.ServiceDetail) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if detail.Grpc.HeaderTransfer == "" {
			return streamer(ctx, desc, cc, method, opts...)
		}

		transfers := strings.Split(detail.Grpc.HeaderTransfer, ",")

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "invalid metadata")
		}

		for _, t := range transfers {
			ops := strings.Split(t, " ")
			if len(ops) < 2 {
				return nil, status.Errorf(codes.InvalidArgument, "invalid header transfer %s", t)
			}
			switch ops[0] {
			case "add":
				if len(ops) != 3 {
					return nil, status.Errorf(codes.InvalidArgument, "invalid header transfer %s", t)
				}
				md.Set(ops[1], ops[2])
			case "edit":
				if len(ops) != 3 {
					return nil, status.Errorf(codes.InvalidArgument, "invalid header transfer %s", t)
				}
				md.Set(ops[1], ops[2])
			case "del":
				if !(2 <= len(ops) && len(ops) <= 3) {
					return nil, status.Errorf(codes.InvalidArgument, "invalid header transfer %s", t)
				}
				md.Delete(ops[1])
			default:
				continue
			}
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
