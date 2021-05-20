// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// BrokerClient is the client API for Broker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrokerClient interface {
	Assignments(ctx context.Context, opts ...grpc.CallOption) (Broker_AssignmentsClient, error)
	Finished(ctx context.Context, in *Assignment, opts ...grpc.CallOption) (*Empty, error)
	Create(ctx context.Context, in *Simulation, opts ...grpc.CallOption) (*SimulationId, error)
	GetSimulation(ctx context.Context, in *SimulationId, opts ...grpc.CallOption) (*Simulation, error)
	GetOppConfig(ctx context.Context, in *SimulationId, opts ...grpc.CallOption) (*OppConfig, error)
	AddTasks(ctx context.Context, in *Tasks, opts ...grpc.CallOption) (*Empty, error)
	SetSource(ctx context.Context, in *Source, opts ...grpc.CallOption) (*Empty, error)
	GetSource(ctx context.Context, in *SimulationId, opts ...grpc.CallOption) (*Source, error)
	AddBinary(ctx context.Context, in *Binary, opts ...grpc.CallOption) (*Empty, error)
	GetBinary(ctx context.Context, in *Arch, opts ...grpc.CallOption) (*Binary, error)
}

type brokerClient struct {
	cc grpc.ClientConnInterface
}

func NewBrokerClient(cc grpc.ClientConnInterface) BrokerClient {
	return &brokerClient{cc}
}

func (c *brokerClient) Assignments(ctx context.Context, opts ...grpc.CallOption) (Broker_AssignmentsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Broker_ServiceDesc.Streams[0], "/service.Broker/Assignments", opts...)
	if err != nil {
		return nil, err
	}
	x := &brokerAssignmentsClient{stream}
	return x, nil
}

type Broker_AssignmentsClient interface {
	Send(*Utilization) error
	Recv() (*Assignment, error)
	grpc.ClientStream
}

type brokerAssignmentsClient struct {
	grpc.ClientStream
}

func (x *brokerAssignmentsClient) Send(m *Utilization) error {
	return x.ClientStream.SendMsg(m)
}

func (x *brokerAssignmentsClient) Recv() (*Assignment, error) {
	m := new(Assignment)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *brokerClient) Finished(ctx context.Context, in *Assignment, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/service.Broker/Finished", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) Create(ctx context.Context, in *Simulation, opts ...grpc.CallOption) (*SimulationId, error) {
	out := new(SimulationId)
	err := c.cc.Invoke(ctx, "/service.Broker/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) GetSimulation(ctx context.Context, in *SimulationId, opts ...grpc.CallOption) (*Simulation, error) {
	out := new(Simulation)
	err := c.cc.Invoke(ctx, "/service.Broker/GetSimulation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) GetOppConfig(ctx context.Context, in *SimulationId, opts ...grpc.CallOption) (*OppConfig, error) {
	out := new(OppConfig)
	err := c.cc.Invoke(ctx, "/service.Broker/GetOppConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) AddTasks(ctx context.Context, in *Tasks, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/service.Broker/AddTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) SetSource(ctx context.Context, in *Source, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/service.Broker/SetSource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) GetSource(ctx context.Context, in *SimulationId, opts ...grpc.CallOption) (*Source, error) {
	out := new(Source)
	err := c.cc.Invoke(ctx, "/service.Broker/GetSource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) AddBinary(ctx context.Context, in *Binary, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/service.Broker/AddBinary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) GetBinary(ctx context.Context, in *Arch, opts ...grpc.CallOption) (*Binary, error) {
	out := new(Binary)
	err := c.cc.Invoke(ctx, "/service.Broker/GetBinary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BrokerServer is the server API for Broker service.
// All implementations must embed UnimplementedBrokerServer
// for forward compatibility
type BrokerServer interface {
	Assignments(Broker_AssignmentsServer) error
	Finished(context.Context, *Assignment) (*Empty, error)
	Create(context.Context, *Simulation) (*SimulationId, error)
	GetSimulation(context.Context, *SimulationId) (*Simulation, error)
	GetOppConfig(context.Context, *SimulationId) (*OppConfig, error)
	AddTasks(context.Context, *Tasks) (*Empty, error)
	SetSource(context.Context, *Source) (*Empty, error)
	GetSource(context.Context, *SimulationId) (*Source, error)
	AddBinary(context.Context, *Binary) (*Empty, error)
	GetBinary(context.Context, *Arch) (*Binary, error)
	mustEmbedUnimplementedBrokerServer()
}

// UnimplementedBrokerServer must be embedded to have forward compatible implementations.
type UnimplementedBrokerServer struct {
}

func (UnimplementedBrokerServer) Assignments(Broker_AssignmentsServer) error {
	return status.Errorf(codes.Unimplemented, "method Assignments not implemented")
}
func (UnimplementedBrokerServer) Finished(context.Context, *Assignment) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Finished not implemented")
}
func (UnimplementedBrokerServer) Create(context.Context, *Simulation) (*SimulationId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedBrokerServer) GetSimulation(context.Context, *SimulationId) (*Simulation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSimulation not implemented")
}
func (UnimplementedBrokerServer) GetOppConfig(context.Context, *SimulationId) (*OppConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOppConfig not implemented")
}
func (UnimplementedBrokerServer) AddTasks(context.Context, *Tasks) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTasks not implemented")
}
func (UnimplementedBrokerServer) SetSource(context.Context, *Source) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetSource not implemented")
}
func (UnimplementedBrokerServer) GetSource(context.Context, *SimulationId) (*Source, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSource not implemented")
}
func (UnimplementedBrokerServer) AddBinary(context.Context, *Binary) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBinary not implemented")
}
func (UnimplementedBrokerServer) GetBinary(context.Context, *Arch) (*Binary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBinary not implemented")
}
func (UnimplementedBrokerServer) mustEmbedUnimplementedBrokerServer() {}

// UnsafeBrokerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrokerServer will
// result in compilation errors.
type UnsafeBrokerServer interface {
	mustEmbedUnimplementedBrokerServer()
}

func RegisterBrokerServer(s grpc.ServiceRegistrar, srv BrokerServer) {
	s.RegisterService(&Broker_ServiceDesc, srv)
}

func _Broker_Assignments_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BrokerServer).Assignments(&brokerAssignmentsServer{stream})
}

type Broker_AssignmentsServer interface {
	Send(*Assignment) error
	Recv() (*Utilization, error)
	grpc.ServerStream
}

type brokerAssignmentsServer struct {
	grpc.ServerStream
}

func (x *brokerAssignmentsServer) Send(m *Assignment) error {
	return x.ServerStream.SendMsg(m)
}

func (x *brokerAssignmentsServer) Recv() (*Utilization, error) {
	m := new(Utilization)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Broker_Finished_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Assignment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).Finished(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Broker/Finished",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).Finished(ctx, req.(*Assignment))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Simulation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Broker/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).Create(ctx, req.(*Simulation))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_GetSimulation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimulationId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).GetSimulation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Broker/GetSimulation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).GetSimulation(ctx, req.(*SimulationId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_GetOppConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimulationId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).GetOppConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Broker/GetOppConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).GetOppConfig(ctx, req.(*SimulationId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_AddTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Tasks)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).AddTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Broker/AddTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).AddTasks(ctx, req.(*Tasks))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_SetSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Source)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).SetSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Broker/SetSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).SetSource(ctx, req.(*Source))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_GetSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimulationId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).GetSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Broker/GetSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).GetSource(ctx, req.(*SimulationId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_AddBinary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Binary)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).AddBinary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Broker/AddBinary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).AddBinary(ctx, req.(*Binary))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_GetBinary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Arch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).GetBinary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Broker/GetBinary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).GetBinary(ctx, req.(*Arch))
	}
	return interceptor(ctx, in, info, handler)
}

// Broker_ServiceDesc is the grpc.ServiceDesc for Broker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Broker",
	HandlerType: (*BrokerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Finished",
			Handler:    _Broker_Finished_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Broker_Create_Handler,
		},
		{
			MethodName: "GetSimulation",
			Handler:    _Broker_GetSimulation_Handler,
		},
		{
			MethodName: "GetOppConfig",
			Handler:    _Broker_GetOppConfig_Handler,
		},
		{
			MethodName: "AddTasks",
			Handler:    _Broker_AddTasks_Handler,
		},
		{
			MethodName: "SetSource",
			Handler:    _Broker_SetSource_Handler,
		},
		{
			MethodName: "GetSource",
			Handler:    _Broker_GetSource_Handler,
		},
		{
			MethodName: "AddBinary",
			Handler:    _Broker_AddBinary_Handler,
		},
		{
			MethodName: "GetBinary",
			Handler:    _Broker_GetBinary_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Assignments",
			Handler:       _Broker_Assignments_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "broker.proto",
}
