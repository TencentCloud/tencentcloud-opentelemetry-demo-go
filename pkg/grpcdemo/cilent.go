package grpcdemo

import (
	"context"
	"fmt"
	tt "go-opentelemetry-demo/pkg/trace"
	"io"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// StartClient 开启客户端
func StartClient() {
	ctx := context.Background()
	tp := tt.InitTraceProvider(ctx)

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":8083", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer func() { _ = conn.Close() }()

	c := NewHelloServiceClient(conn)

	if err := callSayHello(c); err != nil {
		log.Fatal(err)
	}
	if err := callSayHelloClientStream(c); err != nil {
		log.Fatal(err)
	}
	if err := callSayHelloServerStream(c); err != nil {
		log.Fatal(err)
	}
	if err := callSayHelloBidiStream(c); err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)
}

func callSayHello(c HelloServiceClient) error {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-grpcdemo-client-us-east-1",
		"user-id", "some-test-user-id",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	response, err := c.SayHello(ctx, &HelloRequest{Greeting: "World"})
	if err != nil {
		return fmt.Errorf("calling SayHello: %v", err)
	}
	log.Printf("Response from server: %s", response.Reply)
	return nil
}

func callSayHelloClientStream(c HelloServiceClient) error {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-grpcdemo-client-us-east-1",
		"user-id", "some-test-user-id",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.SayHelloClientStream(ctx)
	if err != nil {
		return fmt.Errorf("opening SayHelloClientStream: %v", err)
	}

	for i := 0; i < 5; i++ {
		err = stream.Send(&HelloRequest{Greeting: "World"})

		time.Sleep(time.Duration(i*50) * time.Millisecond)

		if err != nil {
			return fmt.Errorf("sending to SayHelloClientStream: %v", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("closing SayHelloClientStream: %v", err)
	}

	log.Printf("Response from server: %s", response.Reply)
	return nil
}

func callSayHelloServerStream(c HelloServiceClient) error {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-grpcdemo-client-us-east-1",
		"user-id", "some-test-user-id",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.SayHelloServerStream(ctx, &HelloRequest{Greeting: "World"})
	if err != nil {
		return fmt.Errorf("opening SayHelloServerStream: %v", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("receiving from SayHelloServerStream: %v", err)
		}

		log.Printf("Response from server: %s", response.Reply)
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

func callSayHelloBidiStream(c HelloServiceClient) error {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-grpcdemo-client-us-east-1",
		"user-id", "some-test-user-id",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.SayHelloBidiStream(ctx)
	if err != nil {
		return fmt.Errorf("opening SayHelloBidiStream: %v", err)
	}

	serverClosed := make(chan struct{})
	clientClosed := make(chan struct{})

	go func() {
		for i := 0; i < 5; i++ {
			err := stream.Send(&HelloRequest{Greeting: "World"})

			if err != nil {
				// nolint: revive  // This acts as its own main func.
				log.Fatalf("Error when sending to SayHelloBidiStream: %s", err)
			}

			time.Sleep(50 * time.Millisecond)
		}

		err := stream.CloseSend()
		if err != nil {
			// nolint: revive  // This acts as its own main func.
			log.Fatalf("Error when closing SayHelloBidiStream: %s", err)
		}

		clientClosed <- struct{}{}
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				// nolint: revive  // This acts as its own main func.
				log.Fatalf("Error when receiving from SayHelloBidiStream: %s", err)
			}

			log.Printf("Response from server: %s", response.Reply)
			time.Sleep(50 * time.Millisecond)
		}

		serverClosed <- struct{}{}
	}()

	// Wait until client and server both closed the connection.
	<-clientClosed
	<-serverClosed
	return nil
}
