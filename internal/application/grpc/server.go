package grpc

import (
	"astigo/internal/application/grpc/interceptor"
	"astigo/internal/tool"
	"astigo/pkg/proto"
	"google.golang.org/grpc"
)

type GrpcConfig struct {
	Port int `mapstructure:"port"`
}

func NewGrpcServer(fooService proto.FooServiceServer) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryLoggerInterceptor(tool.Logger)),
	)
	server.RegisterService(&proto.FooService_ServiceDesc, fooService)

	return server
}
