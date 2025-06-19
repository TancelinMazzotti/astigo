package cmd

import (
	"astigo/internal/config"
	"astigo/internal/core"
	"astigo/internal/tool"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	serveCmd = &cobra.Command{
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
			server, err := core.NewServer()
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
)
