package main

import (
	"context"
	pb "grpc-demo/proto"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	// SayHello(client)
	// SayHelloAgain(client)
	// SayHelloStream(client)
	SayHelloStreamAll(client)
}

// 一元 RPC
func SayHello(client pb.GreeterClient) {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "duty"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", resp.Message)
}

// 服务端流式 RPC
func SayHelloAgain(client pb.GreeterClient) {
	stream, err := client.SayHelloAgain(context.Background(), &pb.HelloRequest{Name: "duty"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("Greeting: %s", resp.Message)
	}
}

// 客户端流式 RPC
func SayHelloStream(client pb.GreeterClient) {
	stream, err := client.SayHelloStream(context.Background())
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for i := 0; i < 10; i++ {
		if err := stream.Send(&pb.HelloRequest{Name: "duty"}); err != nil {
			log.Fatalf("56 could not greet: %v", err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", resp.Message)
}

// 双向流式 RPC
func SayHelloStreamAll(client pb.GreeterClient) error {
	stream, err := client.SayHelloStreamAll(context.Background())
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return err
	}
	for n := 0; n < 10; n++ {
		err := stream.Send(&pb.HelloRequest{Name: "duty"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not greet: %v", err)
			return err
		}
		log.Printf("Greeting: %s", resp)
	}

	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return err
	}
	return nil
}
