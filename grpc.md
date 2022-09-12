# grpc 的实践

## 一、前期准备

### 1. 安装 protoc

```shell
wget https://github.com/protocolbuffers/protobuf/releases/download/v21.5/protobuf-all-21.5.zip
unzip protobuf-all-21.5.zip && cd protobuf-21.5/
./configure
make
make install
```

安装完成后，可以通过 `protoc --version` 查看是否安装成功。

### 2. 安装 protoc-gen-go

```shell
go get -u github.com/golang/protobuf/protoc-gen-go
go install github.com/golang/protobuf/protoc-gen-go
mv $GOPATH/bin/protoc-gen-go /usr/local/bin/
```

### 3. 初始化一个项目

```shell
mkdir grpc-demo && cd grpc-demo
go mod init grpc-demo
mkdir client server proto
```

## 二、编写 proto 文件

### 1. 创建 proto 文件

```proto
syntax = "proto3";

package helloworld;

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

### 2. 生成 pb.go 文件

```shell
protoc --go_out=plugins=grpc:. ./proto/*.proto
```

> 报错 1：
> go_out: protoc-gen-go: Plugin failed with status code 1.
> 由于没有指定.proto 文件路径，在 proto 文件第三行加上`option go_package ="./proto";`

> 报错 2：
> --go_out: protoc-gen-go: plugins are not supported; use 'protoc --go-grpc_out=...' to generate gRPC
> 由于 proto 版本过高，无法使用`plugins=grpc`，需要使用`--go-grpc_out=.`
> 将生成的 proto-gen-go-grpc 文件移动到/usr/local/bin/目录下
> 使用新的命令：
> protoc --go_out=. --go_opt=paths=source_relative \\
> --go-grpc_out=. --go-grpc_opt=paths=source_relative \\
> ./proto/\*.proto

## 三、编写 GRPC

### 1. 一元 RPC

> 一元 RPC 是指客户端发送一个请求到服务器，并获得一个响应，就像是一个简单的函数调用。

1. Proto

```proto
rpc SayHello (HelloRequest) returns (HelloReply) {};
```

2. server.go

```go
package main

import (
	"context"
	pb "grpc-demo/proto"
	"net"

	"google.golang.org/grpc"
)

type GreeterServer struct{
	pb.UnimplementedGreeterServer
}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + req.Name}, nil
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis,_ := net.Listen("tcp", ":8080")
	server.Serve(lis)
}
```

3. client.go

```go
package main

import (
	"context"
	pb "grpc-demo/proto"
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
	SayHello(client)
}

func SayHello(client pb.GreeterClient) {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "duty"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", resp.Message)
}
```

### 2. 服务端流 RPC

> 服务端流 RPC 是指客户端发送一个请求到服务器，服务器可以返回多个流，直到没有任何消息返回。

1. Protot

```proto
rpc SayHelloAgain (HelloRequest) returns (stream HelloReply) {};
```

2. server.go

```go
package main

import (
	"context"
	pb "grpc-demo/proto"
	"net"

	"google.golang.org/grpc"
)

type GreeterServer struct{
	pb.UnimplementedGreeterServer
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

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis,_ := net.Listen("tcp", ":8080")
	server.Serve(lis)
}

```

3. client.go

```go
package main

import (
	"context"
	pb "grpc-demo/proto"
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
	SayHelloAgain(client)
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

```

### 3. 客户端流 RPC

> 客户端流 RPC 是指客户端发送多个请求到服务器，服务器返回一个响应。

1. Proto

```proto
rpc SayHelloStream (stream HelloRequest) returns (HelloReply) {};
```

2. server.go

```go
package main

import (
	pb "grpc-demo/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
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

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis, _ := net.Listen("tcp", ":8080")
	server.Serve(lis)
}

```

3. client.go

```go
package main

import (
	"context"
	pb "grpc-demo/proto"
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
	SayHelloStream(client)
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

```

### 4. 双向流 RPC

> 双向流 RPC 是指客户端和服务器都可以在双向流中发送多个消息。

1. Proto

```proto
rpc SayHelloStreamAll (stream HelloRequest) returns (stream HelloReply) {};
```

2. server.go

```go
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
```

3. client.go

```go
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
	SayHelloStreamAll(client)
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

```
