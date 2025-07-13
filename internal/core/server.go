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
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

// Config represents the main configuration structure for the application, encompassing all required service settings.
type Config struct {
	Log LoggerConfig `mapstructure:"log"`

	Gin    http2.GinConfig     `mapstructure:"http"`
	Grpc   grpc2.GrpcConfig    `mapstructure:"grpc"`
	Jaeger tracer.JaegerConfig `mapstructure:"jaeger"`
	Auth   struct {
		ClientID string `mapstructure:"client_id"`
		Issuer   string `mapstructure:"issuer"`
	} `mapstructure:"auth"`

	Postgres postgres2.PostgresConfig `mapstructure:"postgres"`
	Nats     nats2.NatsConfig         `mapstructure:"nats"`
	Redis    redis2.RedisConfig       `mapstructure:"redis"`
}

// Server represents the main service structure that holds all essential configurations and dependencies.
type Server struct {
	Config Config
	Logger *zap.Logger

	HttpServer   *http.Server
	GrpcServer   *grpc.Server
	ConsumerNats *event.ConsumerNats
	GinEngine    *gin.Engine

	Provider *oidc.Provider
	Jaeger   *tracer.Jaeger
	Postgres *sql.DB
	Nats     *nats.Conn
	Redis    *redis.Client
}

// Start initializes and runs the HTTP and gRPC servers concurrently and listens for errors and shutdown signals.
func (server *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 2)

	go server.startHTTPServer(errCh)
	go server.startGrpcServer(errCh)
	go server.handleShutdown(ctx, errCh)

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

// startHTTPServer starts an HTTP server using the configuration provided in the Server instance and listens for errors.
func (server *Server) startHTTPServer(errCh chan<- error) {
	http2.StartAt = time.Now()
	server.HttpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", server.Config.Gin.Port),
		Handler: server.GinEngine,
	}

	server.Logger.Info("HTTP server starting")
	if err := server.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errCh <- fmt.Errorf("HTTP server error: %w", err)
	}

}

// startGrpcServer starts the gRPC server using the configured port and listens for errors during its operation.
func (server *Server) startGrpcServer(errCh chan<- error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", server.Config.Grpc.Port))
	if err != nil {
		errCh <- fmt.Errorf("failed to listen for gRPC: %w", err)
		return
	}

	server.Logger.Info("gRPC server starting")
	if err := server.GrpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		errCh <- fmt.Errorf("gRPC server error: %w", err)
	}
}

// handleShutdown manages the server shutdown process by gracefully stopping services and releasing resources.
func (server *Server) handleShutdown(ctx context.Context, errCh chan<- error) {
	<-ctx.Done()
	server.Logger.Info("Shutdown signal received...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP
	server.Logger.Info("HTTP server shutdown...")
	if err := server.HttpServer.Shutdown(shutdownCtx); err != nil {
		server.Logger.Error("HTTP shutdown error", zap.Error(err))
	} else {
		server.Logger.Info("HTTP server shutdown")
	}

	// Shutdown gRPC
	server.Logger.Info("gRPC server shutdown...")
	server.GrpcServer.GracefulStop()
	server.Logger.Info("gRPC server shutdown")

	// Close NatsConsumer
	server.Logger.Info("Nats consumer shutdown...")
	if err := server.ConsumerNats.Close(); err != nil {
		server.Logger.Error("Nats consumer shutdown error", zap.Error(err))
	} else {
		server.Logger.Info("Nats consumer shutdown")
	}

	// Shutdown Postgres
	server.Logger.Info("Postgres shutdown...")
	if err := server.Postgres.Close(); err != nil {
		server.Logger.Error("Postgres shutdown error", zap.Error(err))
	} else {
		server.Logger.Info("Postgres shutdown")
	}

	// Shutdown redis
	server.Logger.Info("Redis shutdown...")
	if err := server.Redis.Close(); err != nil {
		server.Logger.Error("Redis shutdown error", zap.Error(err))
	} else {
		server.Logger.Info("Redis shutdown")
	}

	// Close Nats
	server.Logger.Info("Nats shutdown...")
	server.Nats.Close()
	server.Logger.Info("Nats shutdown")

	// Shutdown Jaeger
	server.Logger.Info("Jaeger shutdown...")
	if err := server.Jaeger.Shutdown(shutdownCtx); err != nil {
		server.Logger.Error("Jaeger shutdown error", zap.Error(err))
	} else {
		server.Logger.Info("Jaeger shutdown")
	}

	server.Logger.Info("Shutdown signal complete")
	errCh <- nil
}

// NewServer initializes a new Server instance with configured dependencies including logging, tracing, and connectors.
// It sets up components such as Jaeger tracer, PostgreSQL, Redis, NATS, OIDC provider, and associated services.
// Returns a fully initialized Server instance or an error if any dependency setup fails.
func NewServer(ctx context.Context, config Config) (*Server, error) {
	var err error
	server := &Server{
		Config: config,
	}

	server.Logger, err = NewLogger(server.Config.Log)
	if err != nil {
		return nil, fmt.Errorf("fail to create logger %w", err)
	}

	server.Logger.Info("create new jaeger tracer")
	server.Jaeger, err = tracer.NewJaeger(ctx, server.Config.Jaeger)
	if err != nil {
		server.Logger.Error("fail to create jaeger tracer", zap.Error(err))
		return nil, fmt.Errorf("fail to create jaeger tracer %w", err)
	}

	server.Logger.Info("create new postgres connector")
	if server.Postgres, err = postgres2.NewPostgres(ctx, server.Config.Postgres); err != nil {
		server.Logger.Error("fail to create postgres connector", zap.Error(err))
		return nil, fmt.Errorf("fail to create postgres connector %w", err)
	}

	server.Logger.Info("create new redis connector")
	if server.Redis, err = redis2.NewRedis(ctx, server.Config.Redis); err != nil {
		server.Logger.Error("fail to create redis connector", zap.Error(err))
		return nil, fmt.Errorf("fail to create nats connector %w", err)
	}

	server.Logger.Info("create new nats connector")
	if server.Nats, err = nats2.NewNats(server.Config.Nats); err != nil {
		server.Logger.Error("fail to create nats connector", zap.Error(err))
		return nil, fmt.Errorf("fail to create nats connector %w", err)
	}

	server.Logger.Info("create new oidc provider")
	server.Provider, err = oidc.NewProvider(context.Background(), server.Config.Auth.Issuer)
	if err != nil {
		server.Logger.Error("fail to create oidc provider", zap.Error(err))
		return nil, fmt.Errorf("failed to create oidc provider: %w", err)
	}

	server.Logger.Info("create new nats consumer")
	server.ConsumerNats, err = event.NewConsumerNats(server.Logger, server.Nats)
	if err != nil {
		server.Logger.Error("fail to create nats consumer", zap.Error(err))
		return nil, fmt.Errorf("fail to create nats consumer %w", err)
	}

	server.Logger.Debug("create new auth services")
	authService := service.NewAuthService(
		server.Logger,
		server.Provider,
		server.Config.Auth.ClientID,
	)

	server.Logger.Debug("create new foo services")
	fooService := service.NewFooService(
		server.Logger,
		postgres2.NewFooPostgres(server.Postgres),
		redis2.NewFooRedis(server.Redis),
		nats2.NewFooNats(server.Nats),
	)

	server.Logger.Debug("create new gin engine")
	server.GinEngine = http2.NewGin(
		server.Config.Gin,
		server.Logger,
		authService,
		http2.NewHealthController(),
		http2.NewFooController(fooService),
	)

	server.Logger.Debug("create new grpc server")
	server.GrpcServer = grpc2.NewGrpcServer(
		server.Logger,
		grpc2.NewFooService(fooService),
	)

	return server, nil
}
