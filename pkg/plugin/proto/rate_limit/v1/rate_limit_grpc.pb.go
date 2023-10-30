// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: rate_limit.proto

package __

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

// RateLimitServiceClient is the client API for RateLimitService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RateLimitServiceClient interface {
	RateLimit(ctx context.Context, in *RateLimitRequest, opts ...grpc.CallOption) (*RateLimitResponse, error)
}

type rateLimitServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRateLimitServiceClient(cc grpc.ClientConnInterface) RateLimitServiceClient {
	return &rateLimitServiceClient{cc}
}

func (c *rateLimitServiceClient) RateLimit(ctx context.Context, in *RateLimitRequest, opts ...grpc.CallOption) (*RateLimitResponse, error) {
	out := new(RateLimitResponse)
	err := c.cc.Invoke(ctx, "/io.opensergo.plugin.proto.rate_limit.RateLimitService/RateLimit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RateLimitServiceServer is the server API for RateLimitService service.
// All implementations must embed UnimplementedRateLimitServiceServer
// for forward compatibility
type RateLimitServiceServer interface {
	RateLimit(context.Context, *RateLimitRequest) (*RateLimitResponse, error)
	mustEmbedUnimplementedRateLimitServiceServer()
}

// UnimplementedRateLimitServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRateLimitServiceServer struct {
}

func (UnimplementedRateLimitServiceServer) RateLimit(context.Context, *RateLimitRequest) (*RateLimitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RateLimit not implemented")
}
func (UnimplementedRateLimitServiceServer) mustEmbedUnimplementedRateLimitServiceServer() {}

// UnsafeRateLimitServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RateLimitServiceServer will
// result in compilation errors.
type UnsafeRateLimitServiceServer interface {
	mustEmbedUnimplementedRateLimitServiceServer()
}

func RegisterRateLimitServiceServer(s grpc.ServiceRegistrar, srv RateLimitServiceServer) {
	s.RegisterService(&RateLimitService_ServiceDesc, srv)
}

func _RateLimitService_RateLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RateLimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimitServiceServer).RateLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/io.opensergo.plugin.proto.rate_limit.RateLimitService/RateLimit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimitServiceServer).RateLimit(ctx, req.(*RateLimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RateLimitService_ServiceDesc is the grpc.ServiceDesc for RateLimitService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RateLimitService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "io.opensergo.plugin.proto.rate_limit.RateLimitService",
	HandlerType: (*RateLimitServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RateLimit",
			Handler:    _RateLimitService_RateLimit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rate_limit.proto",
}