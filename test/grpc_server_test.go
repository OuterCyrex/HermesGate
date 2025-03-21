package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	pb "GoGateway/test/grpc/greet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type greetServer struct {
	pb.UnimplementedGreeterServer
}

func (s *greetServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello, %s!", req.Name)}, nil
}

func TestGrpcServer1(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			fmt.Println("50051")
			return handler(ctx, req)
		},
	))
	pb.RegisterGreeterServer(s, &greetServer{})
	reflection.Register(s)

	log.Printf("gRPC Server started on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func TestGrpcServer2(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			fmt.Println("50052")
			return handler(ctx, req)
		},
	))
	pb.RegisterGreeterServer(s, &greetServer{})
	reflection.Register(s)

	log.Printf("gRPC Server started on port 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
