// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.1
// source: proto/urlshortener.proto

package proto

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
	URLShortener_CreateShort_FullMethodName   = "/urlshortener.URLShortener/CreateShort"
	URLShortener_GetOriginLink_FullMethodName = "/urlshortener.URLShortener/GetOriginLink"
	URLShortener_Shorten_FullMethodName       = "/urlshortener.URLShortener/Shorten"
	URLShortener_GetUserUrls_FullMethodName   = "/urlshortener.URLShortener/GetUserUrls"
	URLShortener_DeleteUrls_FullMethodName    = "/urlshortener.URLShortener/DeleteUrls"
	URLShortener_GetStats_FullMethodName      = "/urlshortener.URLShortener/GetStats"
	URLShortener_Ping_FullMethodName          = "/urlshortener.URLShortener/Ping"
)

// URLShortenerClient is the client API for URLShortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLShortenerClient interface {
	// Создать короткий URL
	CreateShort(ctx context.Context, in *CreateShortRequest, opts ...grpc.CallOption) (*CreateShortResponse, error)
	// Получить оригинальный URL по короткому
	GetOriginLink(ctx context.Context, in *GetOriginLinkRequest, opts ...grpc.CallOption) (*GetOriginLinkResponse, error)
	// Укоротить URL через API
	Shorten(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error)
	// Получить все URL пользователя
	GetUserUrls(ctx context.Context, in *GetUserUrlsRequest, opts ...grpc.CallOption) (*GetUserUrlsResponse, error)
	// Удалить URL пользователя
	DeleteUrls(ctx context.Context, in *DeleteUrlsRequest, opts ...grpc.CallOption) (*DeleteUrlsResponse, error)
	// Получить статистику по URL
	GetStats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error)
	// Проверка доступности сервера
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error)
}

type uRLShortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewURLShortenerClient(cc grpc.ClientConnInterface) URLShortenerClient {
	return &uRLShortenerClient{cc}
}

func (c *uRLShortenerClient) CreateShort(ctx context.Context, in *CreateShortRequest, opts ...grpc.CallOption) (*CreateShortResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateShortResponse)
	err := c.cc.Invoke(ctx, URLShortener_CreateShort_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetOriginLink(ctx context.Context, in *GetOriginLinkRequest, opts ...grpc.CallOption) (*GetOriginLinkResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetOriginLinkResponse)
	err := c.cc.Invoke(ctx, URLShortener_GetOriginLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) Shorten(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenResponse)
	err := c.cc.Invoke(ctx, URLShortener_Shorten_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetUserUrls(ctx context.Context, in *GetUserUrlsRequest, opts ...grpc.CallOption) (*GetUserUrlsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserUrlsResponse)
	err := c.cc.Invoke(ctx, URLShortener_GetUserUrls_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) DeleteUrls(ctx context.Context, in *DeleteUrlsRequest, opts ...grpc.CallOption) (*DeleteUrlsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteUrlsResponse)
	err := c.cc.Invoke(ctx, URLShortener_DeleteUrls_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetStats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, URLShortener_GetStats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PongResponse)
	err := c.cc.Invoke(ctx, URLShortener_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLShortenerServer is the server API for URLShortener service.
// All implementations must embed UnimplementedURLShortenerServer
// for forward compatibility.
type URLShortenerServer interface {
	// Создать короткий URL
	CreateShort(context.Context, *CreateShortRequest) (*CreateShortResponse, error)
	// Получить оригинальный URL по короткому
	GetOriginLink(context.Context, *GetOriginLinkRequest) (*GetOriginLinkResponse, error)
	// Укоротить URL через API
	Shorten(context.Context, *ShortenRequest) (*ShortenResponse, error)
	// Получить все URL пользователя
	GetUserUrls(context.Context, *GetUserUrlsRequest) (*GetUserUrlsResponse, error)
	// Удалить URL пользователя
	DeleteUrls(context.Context, *DeleteUrlsRequest) (*DeleteUrlsResponse, error)
	// Получить статистику по URL
	GetStats(context.Context, *StatsRequest) (*StatsResponse, error)
	// Проверка доступности сервера
	Ping(context.Context, *PingRequest) (*PongResponse, error)
	mustEmbedUnimplementedURLShortenerServer()
}

// UnimplementedURLShortenerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedURLShortenerServer struct{}

func (UnimplementedURLShortenerServer) CreateShort(context.Context, *CreateShortRequest) (*CreateShortResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShort not implemented")
}
func (UnimplementedURLShortenerServer) GetOriginLink(context.Context, *GetOriginLinkRequest) (*GetOriginLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginLink not implemented")
}
func (UnimplementedURLShortenerServer) Shorten(context.Context, *ShortenRequest) (*ShortenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shorten not implemented")
}
func (UnimplementedURLShortenerServer) GetUserUrls(context.Context, *GetUserUrlsRequest) (*GetUserUrlsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserUrls not implemented")
}
func (UnimplementedURLShortenerServer) DeleteUrls(context.Context, *DeleteUrlsRequest) (*DeleteUrlsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUrls not implemented")
}
func (UnimplementedURLShortenerServer) GetStats(context.Context, *StatsRequest) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedURLShortenerServer) Ping(context.Context, *PingRequest) (*PongResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedURLShortenerServer) mustEmbedUnimplementedURLShortenerServer() {}
func (UnimplementedURLShortenerServer) testEmbeddedByValue()                      {}

// UnsafeURLShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLShortenerServer will
// result in compilation errors.
type UnsafeURLShortenerServer interface {
	mustEmbedUnimplementedURLShortenerServer()
}

func RegisterURLShortenerServer(s grpc.ServiceRegistrar, srv URLShortenerServer) {
	// If the following call pancis, it indicates UnimplementedURLShortenerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&URLShortener_ServiceDesc, srv)
}

func _URLShortener_CreateShort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShortRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).CreateShort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_CreateShort_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).CreateShort(ctx, req.(*CreateShortRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetOriginLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOriginLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetOriginLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_GetOriginLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetOriginLink(ctx, req.(*GetOriginLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_Shorten_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).Shorten(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_Shorten_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).Shorten(ctx, req.(*ShortenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetUserUrls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserUrlsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetUserUrls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_GetUserUrls_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetUserUrls(ctx, req.(*GetUserUrlsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_DeleteUrls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUrlsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).DeleteUrls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_DeleteUrls_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).DeleteUrls(ctx, req.(*DeleteUrlsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_GetStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetStats(ctx, req.(*StatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLShortener_ServiceDesc is the grpc.ServiceDesc for URLShortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLShortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "urlshortener.URLShortener",
	HandlerType: (*URLShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateShort",
			Handler:    _URLShortener_CreateShort_Handler,
		},
		{
			MethodName: "GetOriginLink",
			Handler:    _URLShortener_GetOriginLink_Handler,
		},
		{
			MethodName: "Shorten",
			Handler:    _URLShortener_Shorten_Handler,
		},
		{
			MethodName: "GetUserUrls",
			Handler:    _URLShortener_GetUserUrls_Handler,
		},
		{
			MethodName: "DeleteUrls",
			Handler:    _URLShortener_DeleteUrls_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _URLShortener_GetStats_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _URLShortener_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/urlshortener.proto",
}