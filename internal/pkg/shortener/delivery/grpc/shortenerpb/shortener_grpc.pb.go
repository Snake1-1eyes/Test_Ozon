// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: shortener.proto

package shortenerpb

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
	ShortenerService_CreateShortLink_FullMethodName = "/shortener.ShortenerService/CreateShortLink"
	ShortenerService_GetShortLink_FullMethodName    = "/shortener.ShortenerService/GetShortLink"
)

// ShortenerServiceClient is the client API for ShortenerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerServiceClient interface {
	CreateShortLink(ctx context.Context, in *CreateShortLinkRequest, opts ...grpc.CallOption) (*CreateShortLinkResponse, error)
	GetShortLink(ctx context.Context, in *GetShortLinkRequest, opts ...grpc.CallOption) (*GetShortLinkResponse, error)
}

type shortenerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerServiceClient(cc grpc.ClientConnInterface) ShortenerServiceClient {
	return &shortenerServiceClient{cc}
}

func (c *shortenerServiceClient) CreateShortLink(ctx context.Context, in *CreateShortLinkRequest, opts ...grpc.CallOption) (*CreateShortLinkResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateShortLinkResponse)
	err := c.cc.Invoke(ctx, ShortenerService_CreateShortLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerServiceClient) GetShortLink(ctx context.Context, in *GetShortLinkRequest, opts ...grpc.CallOption) (*GetShortLinkResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetShortLinkResponse)
	err := c.cc.Invoke(ctx, ShortenerService_GetShortLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServiceServer is the server API for ShortenerService service.
// All implementations must embed UnimplementedShortenerServiceServer
// for forward compatibility.
type ShortenerServiceServer interface {
	CreateShortLink(context.Context, *CreateShortLinkRequest) (*CreateShortLinkResponse, error)
	GetShortLink(context.Context, *GetShortLinkRequest) (*GetShortLinkResponse, error)
	mustEmbedUnimplementedShortenerServiceServer()
}

// UnimplementedShortenerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedShortenerServiceServer struct{}

func (UnimplementedShortenerServiceServer) CreateShortLink(context.Context, *CreateShortLinkRequest) (*CreateShortLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShortLink not implemented")
}
func (UnimplementedShortenerServiceServer) GetShortLink(context.Context, *GetShortLinkRequest) (*GetShortLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShortLink not implemented")
}
func (UnimplementedShortenerServiceServer) mustEmbedUnimplementedShortenerServiceServer() {}
func (UnimplementedShortenerServiceServer) testEmbeddedByValue()                          {}

// UnsafeShortenerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServiceServer will
// result in compilation errors.
type UnsafeShortenerServiceServer interface {
	mustEmbedUnimplementedShortenerServiceServer()
}

func RegisterShortenerServiceServer(s grpc.ServiceRegistrar, srv ShortenerServiceServer) {
	// If the following call pancis, it indicates UnimplementedShortenerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ShortenerService_ServiceDesc, srv)
}

func _ShortenerService_CreateShortLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShortLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServiceServer).CreateShortLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortenerService_CreateShortLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServiceServer).CreateShortLink(ctx, req.(*CreateShortLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortenerService_GetShortLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetShortLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServiceServer).GetShortLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortenerService_GetShortLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServiceServer).GetShortLink(ctx, req.(*GetShortLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShortenerService_ServiceDesc is the grpc.ServiceDesc for ShortenerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShortenerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.ShortenerService",
	HandlerType: (*ShortenerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateShortLink",
			Handler:    _ShortenerService_CreateShortLink_Handler,
		},
		{
			MethodName: "GetShortLink",
			Handler:    _ShortenerService_GetShortLink_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shortener.proto",
}
