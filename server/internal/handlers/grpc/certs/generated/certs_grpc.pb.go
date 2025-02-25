// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: certs.proto

package certs_service

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

// CertsClient is the client API for Certs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CertsClient interface {
	IssueCert(ctx context.Context, in *CertIssueReq, opts ...grpc.CallOption) (*SignKeys, error)
	RetrieveCert(ctx context.Context, in *CertRetrieveReq, opts ...grpc.CallOption) (*SignPubKey, error)
}

type certsClient struct {
	cc grpc.ClientConnInterface
}

func NewCertsClient(cc grpc.ClientConnInterface) CertsClient {
	return &certsClient{cc}
}

func (c *certsClient) IssueCert(ctx context.Context, in *CertIssueReq, opts ...grpc.CallOption) (*SignKeys, error) {
	out := new(SignKeys)
	err := c.cc.Invoke(ctx, "/Certs/IssueCert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *certsClient) RetrieveCert(ctx context.Context, in *CertRetrieveReq, opts ...grpc.CallOption) (*SignPubKey, error) {
	out := new(SignPubKey)
	err := c.cc.Invoke(ctx, "/Certs/RetrieveCert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CertsServer is the server API for Certs service.
// All implementations must embed UnimplementedCertsServer
// for forward compatibility
type CertsServer interface {
	IssueCert(context.Context, *CertIssueReq) (*SignKeys, error)
	RetrieveCert(context.Context, *CertRetrieveReq) (*SignPubKey, error)
	mustEmbedUnimplementedCertsServer()
}

// UnimplementedCertsServer must be embedded to have forward compatible implementations.
type UnimplementedCertsServer struct {
}

func (UnimplementedCertsServer) IssueCert(context.Context, *CertIssueReq) (*SignKeys, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssueCert not implemented")
}
func (UnimplementedCertsServer) RetrieveCert(context.Context, *CertRetrieveReq) (*SignPubKey, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveCert not implemented")
}
func (UnimplementedCertsServer) mustEmbedUnimplementedCertsServer() {}

// UnsafeCertsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CertsServer will
// result in compilation errors.
type UnsafeCertsServer interface {
	mustEmbedUnimplementedCertsServer()
}

func RegisterCertsServer(s grpc.ServiceRegistrar, srv CertsServer) {
	s.RegisterService(&Certs_ServiceDesc, srv)
}

func _Certs_IssueCert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CertIssueReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertsServer).IssueCert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Certs/IssueCert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertsServer).IssueCert(ctx, req.(*CertIssueReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Certs_RetrieveCert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CertRetrieveReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertsServer).RetrieveCert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Certs/RetrieveCert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertsServer).RetrieveCert(ctx, req.(*CertRetrieveReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Certs_ServiceDesc is the grpc.ServiceDesc for Certs service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Certs_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Certs",
	HandlerType: (*CertsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IssueCert",
			Handler:    _Certs_IssueCert_Handler,
		},
		{
			MethodName: "RetrieveCert",
			Handler:    _Certs_RetrieveCert_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "certs.proto",
}
