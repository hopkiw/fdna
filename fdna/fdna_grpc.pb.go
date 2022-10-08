// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: fdna/fdna.proto

package fdna

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

// FdnaClient is the client API for Fdna service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FdnaClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Gossip(ctx context.Context, in *GossipRequest, opts ...grpc.CallOption) (*GossipResponse, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
}

type fdnaClient struct {
	cc grpc.ClientConnInterface
}

func NewFdnaClient(cc grpc.ClientConnInterface) FdnaClient {
	return &fdnaClient{cc}
}

func (c *fdnaClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/fdna.Fdna/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fdnaClient) Gossip(ctx context.Context, in *GossipRequest, opts ...grpc.CallOption) (*GossipResponse, error) {
	out := new(GossipResponse)
	err := c.cc.Invoke(ctx, "/fdna.Fdna/Gossip", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fdnaClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, "/fdna.Fdna/Heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FdnaServer is the server API for Fdna service.
// All implementations must embed UnimplementedFdnaServer
// for forward compatibility
type FdnaServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Gossip(context.Context, *GossipRequest) (*GossipResponse, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
	mustEmbedUnimplementedFdnaServer()
}

// UnimplementedFdnaServer must be embedded to have forward compatible implementations.
type UnimplementedFdnaServer struct {
}

func (UnimplementedFdnaServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedFdnaServer) Gossip(context.Context, *GossipRequest) (*GossipResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Gossip not implemented")
}
func (UnimplementedFdnaServer) Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedFdnaServer) mustEmbedUnimplementedFdnaServer() {}

// UnsafeFdnaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FdnaServer will
// result in compilation errors.
type UnsafeFdnaServer interface {
	mustEmbedUnimplementedFdnaServer()
}

func RegisterFdnaServer(s grpc.ServiceRegistrar, srv FdnaServer) {
	s.RegisterService(&Fdna_ServiceDesc, srv)
}

func _Fdna_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FdnaServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fdna.Fdna/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FdnaServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fdna_Gossip_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FdnaServer).Gossip(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fdna.Fdna/Gossip",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FdnaServer).Gossip(ctx, req.(*GossipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fdna_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FdnaServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fdna.Fdna/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FdnaServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Fdna_ServiceDesc is the grpc.ServiceDesc for Fdna service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Fdna_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fdna.Fdna",
	HandlerType: (*FdnaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Fdna_Get_Handler,
		},
		{
			MethodName: "Gossip",
			Handler:    _Fdna_Gossip_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _Fdna_Heartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fdna/fdna.proto",
}
