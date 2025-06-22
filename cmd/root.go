package cmd

import (
	"astigo/internal/core"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// Define persistent flags that will be global for your application
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to configuration file")
	rootCmd.PersistentFlags().Int("http.port", 8080, "HTTP server port")
	rootCmd.PersistentFlags().Int("grpc.port", 50051, "gRPC server port")
	rootCmd.PersistentFlags().String("log.level", "info", "logging level (debug, info, warn, error)")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Printf("failed to bind flags: %v\n", err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "astigo",
	Short: "Start API server",
	Long: `Astigo is a versatile API server that supports both HTTP and gRPC protocols. 
It provides a robust configuration system supporting default config, config files, environment variables, and command-line flags.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var config core.Config
		if err := viper.Unmarshal(&config); err != nil {
			return fmt.Errorf("failed to parse configuration: %w", err)
		}

		server, err := core.NewServer(config)
		if err != nil {
			return fmt.Errorf("failed to initialize server: %w", err)
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		if err := server.Start(ctx); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
