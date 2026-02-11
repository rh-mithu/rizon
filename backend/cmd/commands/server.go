package commands

import (
	"context"
	"github.com/rh-mithu/rizon/backend/config"
	"github.com/rh-mithu/rizon/backend/internal/app"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

type Server struct {
	cfg *config.Config
	l   *slog.Logger
}

func ProvideServer(cfg *config.Config, l *slog.Logger) *Server {
	return &Server{
		cfg: cfg,
		l:   l,
	}
}

func (s *Server) RunServerCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "serve",
		Short:   "Run servers",
		Long:    "Run servers",
		Example: "[binary] server [server_name]",
	}

	command.AddCommand(s.httpServerCommand())
	return command
}

func (s *Server) httpServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "http",
		Short:   "http",
		Long:    "Run HTTP server",
		Example: "[binary] server http",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(
				context.Background(),
				syscall.SIGINT,
				syscall.SIGTERM,
			)
			defer stop()

			application, err := app.New(ctx, s.cfg, s.l)
			if err != nil {
				log.Fatalf("Failed to initialize application: %v", err)
			}
			defer application.Close()

			if err := application.Run(); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		},
	}
}
