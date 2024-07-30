// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.2
// source: post.proto

package gen

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
	PostSrv_Index_FullMethodName = "/gen.PostSrv/Index"
)

// PostSrvClient is the client API for PostSrv service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostSrvClient interface {
	Index(ctx context.Context, in *IndexRequest, opts ...grpc.CallOption) (*IndexResponse, error)
}

type postSrvClient struct {
	cc grpc.ClientConnInterface
}

func NewPostSrvClient(cc grpc.ClientConnInterface) PostSrvClient {
	return &postSrvClient{cc}
}

func (c *postSrvClient) Index(ctx context.Context, in *IndexRequest, opts ...grpc.CallOption) (*IndexResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IndexResponse)
	err := c.cc.Invoke(ctx, PostSrv_Index_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostSrvServer is the server API for PostSrv service.
// All implementations must embed UnimplementedPostSrvServer
// for forward compatibility.
type PostSrvServer interface {
	Index(context.Context, *IndexRequest) (*IndexResponse, error)
	mustEmbedUnimplementedPostSrvServer()
}

// UnimplementedPostSrvServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPostSrvServer struct{}

func (UnimplementedPostSrvServer) Index(context.Context, *IndexRequest) (*IndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}
func (UnimplementedPostSrvServer) mustEmbedUnimplementedPostSrvServer() {}
func (UnimplementedPostSrvServer) testEmbeddedByValue()                 {}

// UnsafePostSrvServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostSrvServer will
// result in compilation errors.
type UnsafePostSrvServer interface {
	mustEmbedUnimplementedPostSrvServer()
}

func RegisterPostSrvServer(s grpc.ServiceRegistrar, srv PostSrvServer) {
	// If the following call pancis, it indicates UnimplementedPostSrvServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PostSrv_ServiceDesc, srv)
}

func _PostSrv_Index_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostSrvServer).Index(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostSrv_Index_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostSrvServer).Index(ctx, req.(*IndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PostSrv_ServiceDesc is the grpc.ServiceDesc for PostSrv service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostSrv_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gen.PostSrv",
	HandlerType: (*PostSrvServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Index",
			Handler:    _PostSrv_Index_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "post.proto",
}
