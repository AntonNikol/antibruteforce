// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: bucket_service.proto

package bucketpb

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

// BucketServiceClient is the client API for BucketService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BucketServiceClient interface {
	ResetBucket(ctx context.Context, in *ResetBucketRequest, opts ...grpc.CallOption) (*ResetBucketResponse, error)
}

type bucketServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBucketServiceClient(cc grpc.ClientConnInterface) BucketServiceClient {
	return &bucketServiceClient{cc}
}

func (c *bucketServiceClient) ResetBucket(ctx context.Context, in *ResetBucketRequest, opts ...grpc.CallOption) (*ResetBucketResponse, error) {
	out := new(ResetBucketResponse)
	err := c.cc.Invoke(ctx, "/bucket.BucketService/ResetBucket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BucketServiceServer is the server API for BucketService service.
// All implementations must embed UnimplementedBucketServiceServer
// for forward compatibility
type BucketServiceServer interface {
	ResetBucket(context.Context, *ResetBucketRequest) (*ResetBucketResponse, error)
	mustEmbedUnimplementedBucketServiceServer()
}

// UnimplementedBucketServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBucketServiceServer struct {
}

func (UnimplementedBucketServiceServer) ResetBucket(context.Context, *ResetBucketRequest) (*ResetBucketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetBucket not implemented")
}
func (UnimplementedBucketServiceServer) mustEmbedUnimplementedBucketServiceServer() {}

// UnsafeBucketServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BucketServiceServer will
// result in compilation errors.
type UnsafeBucketServiceServer interface {
	mustEmbedUnimplementedBucketServiceServer()
}

func RegisterBucketServiceServer(s grpc.ServiceRegistrar, srv BucketServiceServer) {
	s.RegisterService(&BucketService_ServiceDesc, srv)
}

func _BucketService_ResetBucket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetBucketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BucketServiceServer).ResetBucket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bucket.BucketService/ResetBucket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BucketServiceServer).ResetBucket(ctx, req.(*ResetBucketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BucketService_ServiceDesc is the grpc.ServiceDesc for BucketService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BucketService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bucket.BucketService",
	HandlerType: (*BucketServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResetBucket",
			Handler:    _BucketService_ResetBucket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bucket_service.proto",
}
