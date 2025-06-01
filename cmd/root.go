package cmd

import (
	"astigo/internal/tool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "astigo",
	Short: "Astigo est une API REST en Go avec Gin",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		tool.Logger.Error("Erreur d'exécution de la commande", zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "Chemin vers le fichier de configuration")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.AddCommand(serveCmd)
}
