package core

import (
	"astigo/internal/application/event"
	grpc2 "astigo/internal/application/grpc"
	http2 "astigo/internal/application/http"
	"astigo/internal/domain/service"
	redis2 "astigo/internal/infrastructure/cache/redis"
	nats2 "astigo/internal/infrastructure/messaging/nats"
	postgres2 "astigo/internal/infrastructure/repository/postgres"
	"astigo/internal/infrastructure/tracer"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Log LoggerConfig `mapstructure:"log"`

	Gin    http2.GinConfig     `mapstructure:"http"`
	Grpc   grpc2.GrpcConfig    `mapstructure:"grpc"`
	Jaeger tracer.JaegerConfig `mapstructure:"jaeger"`

	Postgres postgres2.PostgresConfig `mapstructure:"postgres"`
	Nats     nats2.NatsConfig         `mapstructure:"nats"`
	Redis    redis2.RedisConfig       `mapstructure:"redis"`
}

type Server struct {
	Config Config
	Logger *zap.Logger

	GinEngine    *gin.Engine
	HttpServer   *http.Server
	GrpcServer   *grpc.Server
	ConsumerNats *event.ConsumerNats
	Jaeger       *tracer.Jaeger

	Postgres *sql.DB
	Nats     *nats.Conn
	Redis    *redis.Client
}

func (server *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 2)

	server.HttpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", server.Config.Gin.Port),
		Handler: server.GinEngine,
	}

	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", server.Config.Grpc.Port))
	if err != nil {
		return fmt.Errorf("failed to listen for gRPC: %w", err)
	}

	go server.startHTTPServer(errCh)
	go server.startGrpcServer(grpcLis, errCh)

	go server.handleShutdown(ctx, errCh)

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

func (server *Server) startHTTPServer(errCh chan<- error) {
	http2.StartAt = time.Now()
	server.Logger.Info("HTTP server starting")
	if err := server.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errCh <- fmt.Errorf("http server error: %w", err)
	}

}

func (server *Server) startGrpcServer(lis net.Listener, errCh chan<- error) {
	server.Logger.Info("gRPC server starting...")
	if err := server.GrpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		errCh <- fmt.Errorf("grpc server error: %w", err)
	}
}

func (server *Server) handleShutdown(ctx context.Context, errCh chan<- error) {
	<-ctx.Done()
	server.Logger.Info("Shutdown signal received...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP
	if err := server.HttpServer.Shutdown(shutdownCtx); err != nil {
		server.Logger.Error("HTTP shutdown error", zap.Error(err))
	} else {
		server.Logger.Info("HTTP server stopped")
	}

	// Shutdown gRPC
	server.GrpcServer.GracefulStop()
	server.Logger.Info("gRPC server stopped")

	// Close NatsConsumer
	if err := server.ConsumerNats.Close(); err != nil {
		server.Logger.Error("error while close nats consumer", zap.Error(err))
	} else {
		server.Logger.Info("Nats consumer stopped")
	}

	// Shutdown Postgres
	if err := server.Postgres.Close(); err != nil {
		server.Logger.Error("Postgres shutdown error", zap.Error(err))
	} else {
		server.Logger.Info("Postgres shutdown")
	}

	// Shutdown redis
	if err := server.Redis.Close(); err != nil {
		server.Logger.Error("Redis shutdown error", zap.Error(err))
	} else {
		server.Logger.Info("Redis shutdown")
	}

	// Close Nats
	server.Nats.Close()
	server.Logger.Info("Nats shutdown")

	// Shutdown Jaeger
	if err := server.Jaeger.Shutdown(shutdownCtx); err != nil {
		server.Logger.Error("Jaeger shutdown error", zap.Error(err))
	} else {
		server.Logger.Info("Jaeger shutdown")
	}

	server.Logger.Info("Shutdown complete")
	errCh <- nil
}

func NewServer(config Config) (*Server, error) {
	var err error
	server := &Server{
		Config: config,
	}

	server.Logger, err = NewLogger(server.Config.Log)
	if err != nil {
		return nil, fmt.Errorf("fail to create logger %w", err)
	}

	server.Jaeger, err = tracer.NewJaeger(server.Config.Jaeger)
	if err != nil {
		return nil, fmt.Errorf("fail to create jaeger tracer %w", err)
	}

	if server.Postgres, err = postgres2.NewPostgres(server.Config.Postgres); err != nil {
		return nil, fmt.Errorf("fail to create postgres connector %w", err)
	}

	if server.Redis, err = redis2.NewRedis(server.Config.Redis); err != nil {
		return nil, fmt.Errorf("fail to create nats connector %w", err)
	}

	if server.Nats, err = nats2.NewNats(server.Config.Nats); err != nil {
		return nil, fmt.Errorf("fail to create nats connector %w", err)
	}

	server.ConsumerNats, err = event.NewConsumerNats(server.Logger, server.Nats)
	if err != nil {
		return nil, fmt.Errorf("fail to create nats consumer %w", err)
	}

	fooService := service.NewFooService(
		server.Logger,
		postgres2.NewFooPostgres(server.Postgres),
		redis2.NewFooRedis(server.Redis),
		nats2.NewFooNats(server.Nats),
	)

	server.GinEngine = http2.NewGin(
		server.Config.Gin,
		server.Logger,
		http2.NewHealthController(),
		http2.NewFooController(fooService),
	)

	server.GrpcServer = grpc2.NewGrpcServer(
		server.Logger,
		grpc2.NewFooService(fooService),
	)

	return server, nil
}
