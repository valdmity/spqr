// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: protos/distribution.proto

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
	DistributionService_CreateDistribution_FullMethodName      = "/spqr.DistributionService/CreateDistribution"
	DistributionService_DropDistribution_FullMethodName        = "/spqr.DistributionService/DropDistribution"
	DistributionService_ListDistribution_FullMethodName        = "/spqr.DistributionService/ListDistribution"
	DistributionService_AlterDistributionAttach_FullMethodName = "/spqr.DistributionService/AlterDistributionAttach"
	DistributionService_GetDistribution_FullMethodName         = "/spqr.DistributionService/GetDistribution"
	DistributionService_GetRelationDistribution_FullMethodName = "/spqr.DistributionService/GetRelationDistribution"
)

// DistributionServiceClient is the client API for DistributionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DistributionServiceClient interface {
	CreateDistribution(ctx context.Context, in *CreateDistributionRequest, opts ...grpc.CallOption) (*CreateDistributionReply, error)
	DropDistribution(ctx context.Context, in *DropDistributionRequest, opts ...grpc.CallOption) (*DropDistributionReply, error)
	ListDistribution(ctx context.Context, in *ListDistributionRequest, opts ...grpc.CallOption) (*ListDistributionReply, error)
	AlterDistributionAttach(ctx context.Context, in *AlterDistributionAttachRequest, opts ...grpc.CallOption) (*AlterDistributionAttachReply, error)
	GetDistribution(ctx context.Context, in *GetDistributionRequest, opts ...grpc.CallOption) (*GetDistributionReply, error)
	GetRelationDistribution(ctx context.Context, in *GetRelationDistributionRequest, opts ...grpc.CallOption) (*GetRelationDistributionReply, error)
}

type distributionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDistributionServiceClient(cc grpc.ClientConnInterface) DistributionServiceClient {
	return &distributionServiceClient{cc}
}

func (c *distributionServiceClient) CreateDistribution(ctx context.Context, in *CreateDistributionRequest, opts ...grpc.CallOption) (*CreateDistributionReply, error) {
	out := new(CreateDistributionReply)
	err := c.cc.Invoke(ctx, DistributionService_CreateDistribution_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributionServiceClient) DropDistribution(ctx context.Context, in *DropDistributionRequest, opts ...grpc.CallOption) (*DropDistributionReply, error) {
	out := new(DropDistributionReply)
	err := c.cc.Invoke(ctx, DistributionService_DropDistribution_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributionServiceClient) ListDistribution(ctx context.Context, in *ListDistributionRequest, opts ...grpc.CallOption) (*ListDistributionReply, error) {
	out := new(ListDistributionReply)
	err := c.cc.Invoke(ctx, DistributionService_ListDistribution_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributionServiceClient) AlterDistributionAttach(ctx context.Context, in *AlterDistributionAttachRequest, opts ...grpc.CallOption) (*AlterDistributionAttachReply, error) {
	out := new(AlterDistributionAttachReply)
	err := c.cc.Invoke(ctx, DistributionService_AlterDistributionAttach_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributionServiceClient) GetDistribution(ctx context.Context, in *GetDistributionRequest, opts ...grpc.CallOption) (*GetDistributionReply, error) {
	out := new(GetDistributionReply)
	err := c.cc.Invoke(ctx, DistributionService_GetDistribution_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributionServiceClient) GetRelationDistribution(ctx context.Context, in *GetRelationDistributionRequest, opts ...grpc.CallOption) (*GetRelationDistributionReply, error) {
	out := new(GetRelationDistributionReply)
	err := c.cc.Invoke(ctx, DistributionService_GetRelationDistribution_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DistributionServiceServer is the server API for DistributionService service.
// All implementations must embed UnimplementedDistributionServiceServer
// for forward compatibility
type DistributionServiceServer interface {
	CreateDistribution(context.Context, *CreateDistributionRequest) (*CreateDistributionReply, error)
	DropDistribution(context.Context, *DropDistributionRequest) (*DropDistributionReply, error)
	ListDistribution(context.Context, *ListDistributionRequest) (*ListDistributionReply, error)
	AlterDistributionAttach(context.Context, *AlterDistributionAttachRequest) (*AlterDistributionAttachReply, error)
	GetDistribution(context.Context, *GetDistributionRequest) (*GetDistributionReply, error)
	GetRelationDistribution(context.Context, *GetRelationDistributionRequest) (*GetRelationDistributionReply, error)
	mustEmbedUnimplementedDistributionServiceServer()
}

// UnimplementedDistributionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDistributionServiceServer struct {
}

func (UnimplementedDistributionServiceServer) CreateDistribution(context.Context, *CreateDistributionRequest) (*CreateDistributionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDistribution not implemented")
}
func (UnimplementedDistributionServiceServer) DropDistribution(context.Context, *DropDistributionRequest) (*DropDistributionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DropDistribution not implemented")
}
func (UnimplementedDistributionServiceServer) ListDistribution(context.Context, *ListDistributionRequest) (*ListDistributionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDistribution not implemented")
}
func (UnimplementedDistributionServiceServer) AlterDistributionAttach(context.Context, *AlterDistributionAttachRequest) (*AlterDistributionAttachReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AlterDistributionAttach not implemented")
}
func (UnimplementedDistributionServiceServer) GetDistribution(context.Context, *GetDistributionRequest) (*GetDistributionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDistribution not implemented")
}
func (UnimplementedDistributionServiceServer) GetRelationDistribution(context.Context, *GetRelationDistributionRequest) (*GetRelationDistributionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRelationDistribution not implemented")
}
func (UnimplementedDistributionServiceServer) mustEmbedUnimplementedDistributionServiceServer() {}

// UnsafeDistributionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DistributionServiceServer will
// result in compilation errors.
type UnsafeDistributionServiceServer interface {
	mustEmbedUnimplementedDistributionServiceServer()
}

func RegisterDistributionServiceServer(s grpc.ServiceRegistrar, srv DistributionServiceServer) {
	s.RegisterService(&DistributionService_ServiceDesc, srv)
}

func _DistributionService_CreateDistribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDistributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistributionServiceServer).CreateDistribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DistributionService_CreateDistribution_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistributionServiceServer).CreateDistribution(ctx, req.(*CreateDistributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistributionService_DropDistribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DropDistributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistributionServiceServer).DropDistribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DistributionService_DropDistribution_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistributionServiceServer).DropDistribution(ctx, req.(*DropDistributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistributionService_ListDistribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDistributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistributionServiceServer).ListDistribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DistributionService_ListDistribution_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistributionServiceServer).ListDistribution(ctx, req.(*ListDistributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistributionService_AlterDistributionAttach_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlterDistributionAttachRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistributionServiceServer).AlterDistributionAttach(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DistributionService_AlterDistributionAttach_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistributionServiceServer).AlterDistributionAttach(ctx, req.(*AlterDistributionAttachRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistributionService_GetDistribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDistributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistributionServiceServer).GetDistribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DistributionService_GetDistribution_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistributionServiceServer).GetDistribution(ctx, req.(*GetDistributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistributionService_GetRelationDistribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRelationDistributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistributionServiceServer).GetRelationDistribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DistributionService_GetRelationDistribution_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistributionServiceServer).GetRelationDistribution(ctx, req.(*GetRelationDistributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DistributionService_ServiceDesc is the grpc.ServiceDesc for DistributionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DistributionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spqr.DistributionService",
	HandlerType: (*DistributionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateDistribution",
			Handler:    _DistributionService_CreateDistribution_Handler,
		},
		{
			MethodName: "DropDistribution",
			Handler:    _DistributionService_DropDistribution_Handler,
		},
		{
			MethodName: "ListDistribution",
			Handler:    _DistributionService_ListDistribution_Handler,
		},
		{
			MethodName: "AlterDistributionAttach",
			Handler:    _DistributionService_AlterDistributionAttach_Handler,
		},
		{
			MethodName: "GetDistribution",
			Handler:    _DistributionService_GetDistribution_Handler,
		},
		{
			MethodName: "GetRelationDistribution",
			Handler:    _DistributionService_GetRelationDistribution_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/distribution.proto",
}
