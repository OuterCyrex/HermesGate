package test

import (
	pb "GoGateway/test/grpc/greet"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"testing"
	"time"
)

func BenchmarkGrpcProxyServer(b *testing.B) {
	var wg sync.WaitGroup
	const maxConcurrency = 100
	sem := make(chan struct{}, maxConcurrency)

	for i := 0; i < b.N; i++ {
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			conn, err := grpc.NewClient("127.0.0.1:8099", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewGreeterClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Greeting: %s", r.Message)
		}()
	}

	wg.Wait()
}
