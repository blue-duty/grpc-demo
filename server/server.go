package main

import (
	"context"
	pb "grpc-demo/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

// 一元 RPC
func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + req.Name}, nil
}

// 服务端流式 RPC
func (s *GreeterServer) SayHelloAgain(req *pb.HelloRequest, stream pb.Greeter_SayHelloAgainServer) error {
	for i := 0; i < 10; i++ {
		if err := stream.Send(&pb.HelloReply{Message: "Hello " + req.Name}); err != nil {
			return err
		}
	}
	return nil
}

// 客户端流式 RPC
func (s *GreeterServer) SayHelloStream(stream pb.Greeter_SayHelloStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{Message: "Hello " + "duty"})
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %s", req)
	}
}

// 双向流式 RPC
func (s *GreeterServer) SayHelloStreamAll(stream pb.Greeter_SayHelloStreamAllServer) error {
	n := 0
	for {
		err := stream.Send(&pb.HelloReply{Message: "Hello " + "duty"})
		if err != nil {
			return err
		}
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		log.Printf("resp: %s,%d", req, n)
	}
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis, _ := net.Listen("tcp", ":8080")
	server.Serve(lis)
}
