// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.31.0--rc2
// source: vps.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Vps_ReportStatus_FullMethodName = "/vps.Vps/ReportStatus"
	Vps_SendCommand_FullMethodName  = "/vps.Vps/SendCommand"
)

// VpsClient is the client API for Vps service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VpsClient interface {
	ReportStatus(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[StatusRequest, StatusAck], error)
	SendCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error)
}

type vpsClient struct {
	cc grpc.ClientConnInterface
}

func NewVpsClient(cc grpc.ClientConnInterface) VpsClient {
	return &vpsClient{cc}
}

func (c *vpsClient) ReportStatus(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[StatusRequest, StatusAck], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Vps_ServiceDesc.Streams[0], Vps_ReportStatus_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StatusRequest, StatusAck]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Vps_ReportStatusClient = grpc.ClientStreamingClient[StatusRequest, StatusAck]

func (c *vpsClient) SendCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommandResponse)
	err := c.cc.Invoke(ctx, Vps_SendCommand_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VpsServer is the server API for Vps service.
// All implementations must embed UnimplementedVpsServer
// for forward compatibility.
type VpsServer interface {
	ReportStatus(grpc.ClientStreamingServer[StatusRequest, StatusAck]) error
	SendCommand(context.Context, *CommandRequest) (*CommandResponse, error)
	mustEmbedUnimplementedVpsServer()
}

// UnimplementedVpsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedVpsServer struct{}

func (UnimplementedVpsServer) ReportStatus(grpc.ClientStreamingServer[StatusRequest, StatusAck]) error {
	return status.Errorf(codes.Unimplemented, "method ReportStatus not implemented")
}
func (UnimplementedVpsServer) SendCommand(context.Context, *CommandRequest) (*CommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCommand not implemented")
}
func (UnimplementedVpsServer) mustEmbedUnimplementedVpsServer() {}
func (UnimplementedVpsServer) testEmbeddedByValue()             {}

// UnsafeVpsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VpsServer will
// result in compilation errors.
type UnsafeVpsServer interface {
	mustEmbedUnimplementedVpsServer()
}

func RegisterVpsServer(s grpc.ServiceRegistrar, srv VpsServer) {
	// If the following call pancis, it indicates UnimplementedVpsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Vps_ServiceDesc, srv)
}

func _Vps_ReportStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(VpsServer).ReportStatus(&grpc.GenericServerStream[StatusRequest, StatusAck]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Vps_ReportStatusServer = grpc.ClientStreamingServer[StatusRequest, StatusAck]

func _Vps_SendCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VpsServer).SendCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vps_SendCommand_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VpsServer).SendCommand(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Vps_ServiceDesc is the grpc.ServiceDesc for Vps service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Vps_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vps.Vps",
	HandlerType: (*VpsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendCommand",
			Handler:    _Vps_SendCommand_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReportStatus",
			Handler:       _Vps_ReportStatus_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "vps.proto",
}
