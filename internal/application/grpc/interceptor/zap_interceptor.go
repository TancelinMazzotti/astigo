package interceptor

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"time"
)

func UnaryLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		p, _ := peer.FromContext(ctx)
		clientIP := "unknown"
		if p != nil {
			clientIP = p.Addr.String()
		}

		resp, err := handler(ctx, req)

		statusCode := codes.OK
		if err != nil {
			st, _ := status.FromError(err)
			statusCode = st.Code()
		}

		logger.Info("gRPC request",
			zap.String("status", statusCode.String()),
			zap.String("method", info.FullMethod),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", time.Since(start)),
		)

		return resp, err
	}
}
