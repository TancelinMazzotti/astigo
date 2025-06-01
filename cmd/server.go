package cmd

import (
	"astigo/internal/application/http"
	"astigo/internal/config"
	redis2 "astigo/internal/infrastructure/cache/redis"
	nats2 "astigo/internal/infrastructure/messaging/nats"
	"astigo/internal/infrastructure/repository/postgres"
	"astigo/internal/tool"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Démarre le serveur API",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Charger la configuration avant de démarrer le serveur
		if err := config.Load(); err != nil {
			return fmt.Errorf("la configuration n'a pas été chargée correctement")
		}

		if err := tool.InitLogger(config.Cfg.Log); err != nil {
			return fmt.Errorf("le logger n'a pas été chargée correctement")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialiser l'application
		server, err := NewServer()
		if err != nil {
			return fmt.Errorf("erreur lors de l'initialisation de l'application: %w", err)
		}

		// Lancer le serveur
		if err := server.Start(); err != nil {
			return fmt.Errorf("erreur lors du démarrage de l'application: %w", err)
		}

		return nil
	},
}

type Server struct {
	Postgres  *sql.DB
	Nats      *nats.Conn
	Redis     *redis.Client
	GinRouter *gin.Engine
}

func (server *Server) Start() error {
	server.GinRouter.Run(fmt.Sprintf(":%s", config.Cfg.Gin.Port))
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

	server.GinRouter = http.NewGin()

	return server, nil
}
