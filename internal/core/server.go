package core

import (
	grpc2 "astigo/internal/application/grpc"
	http2 "astigo/internal/application/http"
	"astigo/internal/config"
	"astigo/internal/domain/service"
	redis2 "astigo/internal/infrastructure/cache/redis"
	nats2 "astigo/internal/infrastructure/messaging/nats"
	postgres2 "astigo/internal/infrastructure/repository/postgres"
	"astigo/internal/tool"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Postgres   *sql.DB
	Nats       *nats.Conn
	Redis      *redis.Client
	GinEngine  *gin.Engine
	GrpcServer *grpc.Server
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

	// Attend une erreur ou un arrêt normal
	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

func (server *Server) startHTTPServer(errCh chan<- error) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Cfg.Gin.Port),
		Handler: server.GinEngine,
	}

	go func() {
		http2.StartAt = time.Now()
		tool.Logger.Info("HTTP server starting")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server error: %w", err)
		}
	}()

	return srv
}

func (server *Server) startGrpcServer(lis net.Listener, errCh chan<- error) {
	go func() {
		tool.Logger.Info("gRPC server starting...")
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

		errCh <- nil // signal volontaire
	}()
}

func NewServer() (*Server, error) {
	var err error
	server := &Server{}

	if server.Postgres, err = postgres2.NewPostgres(config.Cfg.Postgres); err != nil {
		return nil, fmt.Errorf("fail to create postgres connector %w", err)
	}

	if server.Redis, err = redis2.NewRedis(config.Cfg.Redis); err != nil {
		return nil, fmt.Errorf("fail to create nats connector %w", err)
	}

	if server.Nats, err = nats2.NewNats(config.Cfg.Nats); err != nil {
		return nil, fmt.Errorf("fail to create nats connector %w", err)
	}

	fooService := service.NewFooService(
		postgres2.NewFooPostgres(server.Postgres),
		redis2.NewFooRedis(server.Redis),
		nats2.NewFooNats(server.Nats),
	)

	server.GinEngine = http2.NewGin(
		http2.NewHealthController(),
		http2.NewFooController(fooService),
	)

	server.GrpcServer = grpc2.NewGrpcServer(
		grpc2.NewFooService(fooService),
	)

	return server, nil
}
