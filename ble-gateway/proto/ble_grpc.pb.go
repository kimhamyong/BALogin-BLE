// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/ble.proto

package ble

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

// BLEServiceClient is the client API for BLEService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BLEServiceClient interface {
	SendDeviceStatus(ctx context.Context, in *DeviceStatus, opts ...grpc.CallOption) (*Response, error)
}

type bLEServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBLEServiceClient(cc grpc.ClientConnInterface) BLEServiceClient {
	return &bLEServiceClient{cc}
}

func (c *bLEServiceClient) SendDeviceStatus(ctx context.Context, in *DeviceStatus, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/ble.BLEService/SendDeviceStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BLEServiceServer is the server API for BLEService service.
// All implementations must embed UnimplementedBLEServiceServer
// for forward compatibility
type BLEServiceServer interface {
	SendDeviceStatus(context.Context, *DeviceStatus) (*Response, error)
	mustEmbedUnimplementedBLEServiceServer()
}

// UnimplementedBLEServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBLEServiceServer struct {
}

func (UnimplementedBLEServiceServer) SendDeviceStatus(context.Context, *DeviceStatus) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendDeviceStatus not implemented")
}
func (UnimplementedBLEServiceServer) mustEmbedUnimplementedBLEServiceServer() {}

// UnsafeBLEServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BLEServiceServer will
// result in compilation errors.
type UnsafeBLEServiceServer interface {
	mustEmbedUnimplementedBLEServiceServer()
}

func RegisterBLEServiceServer(s grpc.ServiceRegistrar, srv BLEServiceServer) {
	s.RegisterService(&BLEService_ServiceDesc, srv)
}

func _BLEService_SendDeviceStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeviceStatus)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BLEServiceServer).SendDeviceStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ble.BLEService/SendDeviceStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BLEServiceServer).SendDeviceStatus(ctx, req.(*DeviceStatus))
	}
	return interceptor(ctx, in, info, handler)
}

// BLEService_ServiceDesc is the grpc.ServiceDesc for BLEService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BLEService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ble.BLEService",
	HandlerType: (*BLEServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendDeviceStatus",
			Handler:    _BLEService_SendDeviceStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ble.proto",
}
