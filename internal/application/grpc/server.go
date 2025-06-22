package grpc

import (
	"astigo/internal/application/grpc/interceptor"
	"astigo/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcConfig struct {
	Port int `mapstructure:"port"`
}

func NewGrpcServer(logger *zap.Logger, fooService proto.FooServiceServer) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryLoggerInterceptor(logger)),
	)
	server.RegisterService(&proto.FooService_ServiceDesc, fooService)

	return server
}
