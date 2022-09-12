// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: proto/helloworld.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GreeterClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
	// Server streaming
	SayHelloAgain(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (Greeter_SayHelloAgainClient, error)
	// Client streaming
	SayHelloStream(ctx context.Context, opts ...grpc.CallOption) (Greeter_SayHelloStreamClient, error)
	// Bidirectional streaming
	SayHelloStreamAll(ctx context.Context, opts ...grpc.CallOption) (Greeter_SayHelloStreamAllClient, error)
}

type greeterClient struct {
	cc grpc.ClientConnInterface
}

func NewGreeterClient(cc grpc.ClientConnInterface) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) SayHelloAgain(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (Greeter_SayHelloAgainClient, error) {
	stream, err := c.cc.NewStream(ctx, &Greeter_ServiceDesc.Streams[0], "/helloworld.Greeter/SayHelloAgain", opts...)
	if err != nil {
		return nil, err
	}
	x := &greeterSayHelloAgainClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Greeter_SayHelloAgainClient interface {
	Recv() (*HelloReply, error)
	grpc.ClientStream
}

type greeterSayHelloAgainClient struct {
	grpc.ClientStream
}

func (x *greeterSayHelloAgainClient) Recv() (*HelloReply, error) {
	m := new(HelloReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *greeterClient) SayHelloStream(ctx context.Context, opts ...grpc.CallOption) (Greeter_SayHelloStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Greeter_ServiceDesc.Streams[1], "/helloworld.Greeter/SayHelloStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &greeterSayHelloStreamClient{stream}
	return x, nil
}

type Greeter_SayHelloStreamClient interface {
	Send(*HelloRequest) error
	CloseAndRecv() (*HelloReply, error)
	grpc.ClientStream
}

type greeterSayHelloStreamClient struct {
	grpc.ClientStream
}

func (x *greeterSayHelloStreamClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *greeterSayHelloStreamClient) CloseAndRecv() (*HelloReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(HelloReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *greeterClient) SayHelloStreamAll(ctx context.Context, opts ...grpc.CallOption) (Greeter_SayHelloStreamAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &Greeter_ServiceDesc.Streams[2], "/helloworld.Greeter/SayHelloStreamAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &greeterSayHelloStreamAllClient{stream}
	return x, nil
}

type Greeter_SayHelloStreamAllClient interface {
	Send(*HelloRequest) error
	Recv() (*HelloReply, error)
	grpc.ClientStream
}

type greeterSayHelloStreamAllClient struct {
	grpc.ClientStream
}

func (x *greeterSayHelloStreamAllClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *greeterSayHelloStreamAllClient) Recv() (*HelloReply, error) {
	m := new(HelloReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GreeterServer is the server API for Greeter service.
// All implementations must embed UnimplementedGreeterServer
// for forward compatibility
type GreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	// Server streaming
	SayHelloAgain(*HelloRequest, Greeter_SayHelloAgainServer) error
	// Client streaming
	SayHelloStream(Greeter_SayHelloStreamServer) error
	// Bidirectional streaming
	SayHelloStreamAll(Greeter_SayHelloStreamAllServer) error
	mustEmbedUnimplementedGreeterServer()
}

// UnimplementedGreeterServer must be embedded to have forward compatible implementations.
type UnimplementedGreeterServer struct {
}

func (UnimplementedGreeterServer) SayHello(context.Context, *HelloRequest) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedGreeterServer) SayHelloAgain(*HelloRequest, Greeter_SayHelloAgainServer) error {
	return status.Errorf(codes.Unimplemented, "method SayHelloAgain not implemented")
}
func (UnimplementedGreeterServer) SayHelloStream(Greeter_SayHelloStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SayHelloStream not implemented")
}
func (UnimplementedGreeterServer) SayHelloStreamAll(Greeter_SayHelloStreamAllServer) error {
	return status.Errorf(codes.Unimplemented, "method SayHelloStreamAll not implemented")
}
func (UnimplementedGreeterServer) mustEmbedUnimplementedGreeterServer() {}

// UnsafeGreeterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GreeterServer will
// result in compilation errors.
type UnsafeGreeterServer interface {
	mustEmbedUnimplementedGreeterServer()
}

func RegisterGreeterServer(s grpc.ServiceRegistrar, srv GreeterServer) {
	s.RegisterService(&Greeter_ServiceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_SayHelloAgain_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(HelloRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GreeterServer).SayHelloAgain(m, &greeterSayHelloAgainServer{stream})
}

type Greeter_SayHelloAgainServer interface {
	Send(*HelloReply) error
	grpc.ServerStream
}

type greeterSayHelloAgainServer struct {
	grpc.ServerStream
}

func (x *greeterSayHelloAgainServer) Send(m *HelloReply) error {
	return x.ServerStream.SendMsg(m)
}

func _Greeter_SayHelloStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GreeterServer).SayHelloStream(&greeterSayHelloStreamServer{stream})
}

type Greeter_SayHelloStreamServer interface {
	SendAndClose(*HelloReply) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type greeterSayHelloStreamServer struct {
	grpc.ServerStream
}

func (x *greeterSayHelloStreamServer) SendAndClose(m *HelloReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *greeterSayHelloStreamServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Greeter_SayHelloStreamAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GreeterServer).SayHelloStreamAll(&greeterSayHelloStreamAllServer{stream})
}

type Greeter_SayHelloStreamAllServer interface {
	Send(*HelloReply) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type greeterSayHelloStreamAllServer struct {
	grpc.ServerStream
}

func (x *greeterSayHelloStreamAllServer) Send(m *HelloReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *greeterSayHelloStreamAllServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Greeter_ServiceDesc is the grpc.ServiceDesc for Greeter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Greeter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SayHelloAgain",
			Handler:       _Greeter_SayHelloAgain_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SayHelloStream",
			Handler:       _Greeter_SayHelloStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "SayHelloStreamAll",
			Handler:       _Greeter_SayHelloStreamAll_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/helloworld.proto",
}
