package main

import (
	"github.com/rh-mithu/rizon/backend/cmd/commands"
	"github.com/rh-mithu/rizon/backend/config"
	"github.com/rh-mithu/rizon/backend/pkg/logger"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "analytics",
		Short: "Shikho Analytics management CLI tool",
	}
	err := godotenv.Load()
	if err != nil {
		slog.Info("Error loading .env file")
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	l := logger.NewSlog(cfg)
	err = os.Setenv("TZ", "UTC")
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}
	cmd := commands.ProvideServer(cfg, l)
	rootCmd.AddCommand(cmd.RunServerCommand())
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Command executed with error", "cause", err)
		os.Exit(1)
	}
}
