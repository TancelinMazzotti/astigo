package core

import (
	grpc2 "astigo/internal/application/grpc"
	http2 "astigo/internal/application/http"
	"astigo/internal/domain/service"
	redis2 "astigo/internal/infrastructure/cache/redis"
	nats2 "astigo/internal/infrastructure/messaging/nats"
	postgres2 "astigo/internal/infrastructure/repository/postgres"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Log LoggerConfig `mapstructure:"log"`

	Gin  http2.GinConfig  `mapstructure:"http"`
	Grpc grpc2.GrpcConfig `mapstructure:"grpc"`

	Postgres postgres2.PostgresConfig `mapstructure:"postgres"`
	Nats     nats2.NatsConfig         `mapstructure:"nats"`
	Redis    redis2.RedisConfig       `mapstructure:"redis"`
}

type Server struct {
	Config Config
	Logger *zap.Logger

	GinEngine  *gin.Engine
	GrpcServer *grpc.Server

	Postgres *sql.DB
	Nats     *nats.Conn
	Redis    *redis.Client
}

func (server *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 2)

	httpSrv := server.startHTTPServer(errCh)
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		return fmt.Errorf("failed to listen for gRPC: %w", err)
	}

	server.startGrpcServer(grpcLis, errCh)
	server.handleShutdown(ctx, httpSrv, errCh)

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

func (server *Server) startHTTPServer(errCh chan<- error) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.Config.Gin.Port),
		Handler: server.GinEngine,
	}

	go func() {
		http2.StartAt = time.Now()
		server.Logger.Info("HTTP server starting")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server error: %w", err)
		}
	}()

	return srv
}

func (server *Server) startGrpcServer(lis net.Listener, errCh chan<- error) {
	go func() {
		server.Logger.Info("gRPC server starting...")
		if err := server.GrpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errCh <- fmt.Errorf("grpc server error: %w", err)
		}
	}()
}

func (server *Server) handleShutdown(ctx context.Context, httpSrv *http.Server, errCh chan<- error) {
	go func() {
		<-ctx.Done()
		log.Println("Shutdown signal received...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown HTTP
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP shutdown error: %v", err)
		} else {
			log.Println("HTTP server stopped")
		}

		// Shutdown gRPC
		server.GrpcServer.GracefulStop()
		log.Println("gRPC server stopped")

		errCh <- nil
	}()
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

	if server.Postgres, err = postgres2.NewPostgres(server.Config.Postgres); err != nil {
		return nil, fmt.Errorf("fail to create postgres connector %w", err)
	}

	if server.Redis, err = redis2.NewRedis(server.Config.Redis); err != nil {
		return nil, fmt.Errorf("fail to create nats connector %w", err)
	}

	if server.Nats, err = nats2.NewNats(server.Config.Nats); err != nil {
		return nil, fmt.Errorf("fail to create nats connector %w", err)
	}

	fooService := service.NewFooService(
		server.Logger,
		postgres2.NewFooPostgres(server.Postgres),
		redis2.NewFooRedis(server.Redis),
		nats2.NewFooNats(server.Nats),
	)

	server.GinEngine = http2.NewGin(
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
