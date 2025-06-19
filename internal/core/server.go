package core

import (
	"astigo/internal/application/http"
	"astigo/internal/config"
	"astigo/internal/domain/service"
	redis2 "astigo/internal/infrastructure/cache/redis"
	nats2 "astigo/internal/infrastructure/messaging/nats"
	"astigo/internal/infrastructure/repository/postgres"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"time"
)

var StartAt time.Time

type Server struct {
	Postgres   *sql.DB
	Nats       *nats.Conn
	Redis      *redis.Client
	GinEngine  *gin.Engine
	GrpcServer *grpc.Server
}

func (server *Server) Start() error {
	StartAt = time.Now()

	http.StartAt = StartAt
	if err := server.GinEngine.Run(fmt.Sprintf(":%s", config.Cfg.Gin.Port)); err != nil {
		return fmt.Errorf("fail ro run gin server %w", err)
	}

	return nil
}

func NewServer() (*Server, error) {
	var err error
	server := &Server{}

	if server.Postgres, err = postgres.NewPostgres(config.Cfg.Postgres); err != nil {
		return nil, fmt.Errorf("fail to create postgres connector %w", err)
	}

	if server.Redis, err = redis2.NewRedis(config.Cfg.Redis); err != nil {
		return nil, fmt.Errorf("fail to create nats connector %w", err)
	}

	if server.Nats, err = nats2.NewNats(config.Cfg.Nats); err != nil {
		return nil, fmt.Errorf("fail to create nats connector %w", err)
	}

	fooService := service.NewFooService(
		postgres.NewFooPostgres(server.Postgres),
		redis2.NewFooRedis(server.Redis),
		nats2.NewFooNats(server.Nats),
	)

	server.GinEngine = http.NewGin(
		http.NewHealthController(),
		http.NewFooController(fooService),
	)

	return server, nil
}
