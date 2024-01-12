// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: protos/balancer.proto

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

const (
	BalancerService_ReloadRequired_FullMethodName = "/spqr.BalancerService/ReloadRequired"
)

// BalancerServiceClient is the client API for BalancerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BalancerServiceClient interface {
	// Reload of balancer required when configuration of shards is changed, when new range was added, or when balancer not initialized yet
	ReloadRequired(ctx context.Context, in *ReloadRequest, opts ...grpc.CallOption) (*ReloadReply, error)
}

type balancerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBalancerServiceClient(cc grpc.ClientConnInterface) BalancerServiceClient {
	return &balancerServiceClient{cc}
}

func (c *balancerServiceClient) ReloadRequired(ctx context.Context, in *ReloadRequest, opts ...grpc.CallOption) (*ReloadReply, error) {
	out := new(ReloadReply)
	err := c.cc.Invoke(ctx, BalancerService_ReloadRequired_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BalancerServiceServer is the server API for BalancerService service.
// All implementations must embed UnimplementedBalancerServiceServer
// for forward compatibility
type BalancerServiceServer interface {
	// Reload of balancer required when configuration of shards is changed, when new range was added, or when balancer not initialized yet
	ReloadRequired(context.Context, *ReloadRequest) (*ReloadReply, error)
	mustEmbedUnimplementedBalancerServiceServer()
}

// UnimplementedBalancerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBalancerServiceServer struct {
}

func (UnimplementedBalancerServiceServer) ReloadRequired(context.Context, *ReloadRequest) (*ReloadReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReloadRequired not implemented")
}
func (UnimplementedBalancerServiceServer) mustEmbedUnimplementedBalancerServiceServer() {}

// UnsafeBalancerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BalancerServiceServer will
// result in compilation errors.
type UnsafeBalancerServiceServer interface {
	mustEmbedUnimplementedBalancerServiceServer()
}

func RegisterBalancerServiceServer(s grpc.ServiceRegistrar, srv BalancerServiceServer) {
	s.RegisterService(&BalancerService_ServiceDesc, srv)
}

func _BalancerService_ReloadRequired_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServiceServer).ReloadRequired(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BalancerService_ReloadRequired_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServiceServer).ReloadRequired(ctx, req.(*ReloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BalancerService_ServiceDesc is the grpc.ServiceDesc for BalancerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BalancerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spqr.BalancerService",
	HandlerType: (*BalancerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReloadRequired",
			Handler:    _BalancerService_ReloadRequired_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/balancer.proto",
}
