// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: userfav.proto

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

// UserFavClient is the client API for UserFav service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserFavClient interface {
	GetFavList(ctx context.Context, in *UserFavRequest, opts ...grpc.CallOption) (*UserFavListResponse, error)
	AddUserFav(ctx context.Context, in *UserFavRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteUserFav(ctx context.Context, in *UserFavRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetUserFavDetail(ctx context.Context, in *UserFavRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type userFavClient struct {
	cc grpc.ClientConnInterface
}

func NewUserFavClient(cc grpc.ClientConnInterface) UserFavClient {
	return &userFavClient{cc}
}

func (c *userFavClient) GetFavList(ctx context.Context, in *UserFavRequest, opts ...grpc.CallOption) (*UserFavListResponse, error) {
	out := new(UserFavListResponse)
	err := c.cc.Invoke(ctx, "/UserFav/GetFavList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userFavClient) AddUserFav(ctx context.Context, in *UserFavRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/UserFav/AddUserFav", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userFavClient) DeleteUserFav(ctx context.Context, in *UserFavRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/UserFav/DeleteUserFav", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userFavClient) GetUserFavDetail(ctx context.Context, in *UserFavRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/UserFav/GetUserFavDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserFavServer is the server API for UserFav service.
// All implementations must embed UnimplementedUserFavServer
// for forward compatibility
type UserFavServer interface {
	GetFavList(context.Context, *UserFavRequest) (*UserFavListResponse, error)
	AddUserFav(context.Context, *UserFavRequest) (*emptypb.Empty, error)
	DeleteUserFav(context.Context, *UserFavRequest) (*emptypb.Empty, error)
	GetUserFavDetail(context.Context, *UserFavRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedUserFavServer()
}

// UnimplementedUserFavServer must be embedded to have forward compatible implementations.
type UnimplementedUserFavServer struct {
}

func (UnimplementedUserFavServer) GetFavList(context.Context, *UserFavRequest) (*UserFavListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFavList not implemented")
}
func (UnimplementedUserFavServer) AddUserFav(context.Context, *UserFavRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserFav not implemented")
}
func (UnimplementedUserFavServer) DeleteUserFav(context.Context, *UserFavRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserFav not implemented")
}
func (UnimplementedUserFavServer) GetUserFavDetail(context.Context, *UserFavRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserFavDetail not implemented")
}
func (UnimplementedUserFavServer) mustEmbedUnimplementedUserFavServer() {}

// UnsafeUserFavServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserFavServer will
// result in compilation errors.
type UnsafeUserFavServer interface {
	mustEmbedUnimplementedUserFavServer()
}

func RegisterUserFavServer(s grpc.ServiceRegistrar, srv UserFavServer) {
	s.RegisterService(&UserFav_ServiceDesc, srv)
}

func _UserFav_GetFavList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserFavRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserFavServer).GetFavList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UserFav/GetFavList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserFavServer).GetFavList(ctx, req.(*UserFavRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserFav_AddUserFav_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserFavRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserFavServer).AddUserFav(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UserFav/AddUserFav",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserFavServer).AddUserFav(ctx, req.(*UserFavRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserFav_DeleteUserFav_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserFavRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserFavServer).DeleteUserFav(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UserFav/DeleteUserFav",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserFavServer).DeleteUserFav(ctx, req.(*UserFavRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserFav_GetUserFavDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserFavRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserFavServer).GetUserFavDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UserFav/GetUserFavDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserFavServer).GetUserFavDetail(ctx, req.(*UserFavRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserFav_ServiceDesc is the grpc.ServiceDesc for UserFav service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserFav_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "UserFav",
	HandlerType: (*UserFavServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFavList",
			Handler:    _UserFav_GetFavList_Handler,
		},
		{
			MethodName: "AddUserFav",
			Handler:    _UserFav_AddUserFav_Handler,
		},
		{
			MethodName: "DeleteUserFav",
			Handler:    _UserFav_DeleteUserFav_Handler,
		},
		{
			MethodName: "GetUserFavDetail",
			Handler:    _UserFav_GetUserFavDetail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "userfav.proto",
}
