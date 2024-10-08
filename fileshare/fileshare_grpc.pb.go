// Copyright 2015 gRPC authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: fileshare/fileshare.proto

package fileshare

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
	FileshareService_HasFile_FullMethodName = "/fileshare.FileshareService/HasFile"
)

// FileshareServiceClient is the client API for FileshareService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileshareServiceClient interface {
	HasFile(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResponse, error)
}

type fileshareServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileshareServiceClient(cc grpc.ClientConnInterface) FileshareServiceClient {
	return &fileshareServiceClient{cc}
}

func (c *fileshareServiceClient) HasFile(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, FileshareService_HasFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileshareServiceServer is the server API for FileshareService service.
// All implementations must embed UnimplementedFileshareServiceServer
// for forward compatibility.
type FileshareServiceServer interface {
	HasFile(context.Context, *MessageRequest) (*MessageResponse, error)
	mustEmbedUnimplementedFileshareServiceServer()
}

// UnimplementedFileshareServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFileshareServiceServer struct{}

func (UnimplementedFileshareServiceServer) HasFile(context.Context, *MessageRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasFile not implemented")
}
func (UnimplementedFileshareServiceServer) mustEmbedUnimplementedFileshareServiceServer() {}
func (UnimplementedFileshareServiceServer) testEmbeddedByValue()                          {}

// UnsafeFileshareServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileshareServiceServer will
// result in compilation errors.
type UnsafeFileshareServiceServer interface {
	mustEmbedUnimplementedFileshareServiceServer()
}

func RegisterFileshareServiceServer(s grpc.ServiceRegistrar, srv FileshareServiceServer) {
	// If the following call pancis, it indicates UnimplementedFileshareServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FileshareService_ServiceDesc, srv)
}

func _FileshareService_HasFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileshareServiceServer).HasFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileshareService_HasFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileshareServiceServer).HasFile(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FileshareService_ServiceDesc is the grpc.ServiceDesc for FileshareService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileshareService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fileshare.FileshareService",
	HandlerType: (*FileshareServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HasFile",
			Handler:    _FileshareService_HasFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fileshare/fileshare.proto",
}
