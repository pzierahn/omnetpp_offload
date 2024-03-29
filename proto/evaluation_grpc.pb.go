// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: evaluation.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EvaluationClient is the client API for Evaluation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EvaluationClient interface {
	ClockSync(ctx context.Context, in *Clock, opts ...grpc.CallOption) (*Clock, error)
	Logs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Evaluation_LogsClient, error)
}

type evaluationClient struct {
	cc grpc.ClientConnInterface
}

func NewEvaluationClient(cc grpc.ClientConnInterface) EvaluationClient {
	return &evaluationClient{cc}
}

func (c *evaluationClient) ClockSync(ctx context.Context, in *Clock, opts ...grpc.CallOption) (*Clock, error) {
	out := new(Clock)
	err := c.cc.Invoke(ctx, "/service.Evaluation/ClockSync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evaluationClient) Logs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Evaluation_LogsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Evaluation_ServiceDesc.Streams[0], "/service.Evaluation/Logs", opts...)
	if err != nil {
		return nil, err
	}
	x := &evaluationLogsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Evaluation_LogsClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type evaluationLogsClient struct {
	grpc.ClientStream
}

func (x *evaluationLogsClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EvaluationServer is the server API for Evaluation service.
// All implementations must embed UnimplementedEvaluationServer
// for forward compatibility
type EvaluationServer interface {
	ClockSync(context.Context, *Clock) (*Clock, error)
	Logs(*emptypb.Empty, Evaluation_LogsServer) error
	mustEmbedUnimplementedEvaluationServer()
}

// UnimplementedEvaluationServer must be embedded to have forward compatible implementations.
type UnimplementedEvaluationServer struct {
}

func (UnimplementedEvaluationServer) ClockSync(context.Context, *Clock) (*Clock, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClockSync not implemented")
}
func (UnimplementedEvaluationServer) Logs(*emptypb.Empty, Evaluation_LogsServer) error {
	return status.Errorf(codes.Unimplemented, "method Logs not implemented")
}
func (UnimplementedEvaluationServer) mustEmbedUnimplementedEvaluationServer() {}

// UnsafeEvaluationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EvaluationServer will
// result in compilation errors.
type UnsafeEvaluationServer interface {
	mustEmbedUnimplementedEvaluationServer()
}

func RegisterEvaluationServer(s grpc.ServiceRegistrar, srv EvaluationServer) {
	s.RegisterService(&Evaluation_ServiceDesc, srv)
}

func _Evaluation_ClockSync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Clock)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvaluationServer).ClockSync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Evaluation/ClockSync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvaluationServer).ClockSync(ctx, req.(*Clock))
	}
	return interceptor(ctx, in, info, handler)
}

func _Evaluation_Logs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EvaluationServer).Logs(m, &evaluationLogsServer{stream})
}

type Evaluation_LogsServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type evaluationLogsServer struct {
	grpc.ServerStream
}

func (x *evaluationLogsServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

// Evaluation_ServiceDesc is the grpc.ServiceDesc for Evaluation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Evaluation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Evaluation",
	HandlerType: (*EvaluationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClockSync",
			Handler:    _Evaluation_ClockSync_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Logs",
			Handler:       _Evaluation_Logs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "evaluation.proto",
}
