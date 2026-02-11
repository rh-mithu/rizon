package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/rh-mithu/rizon/backend/config"
	"github.com/rh-mithu/rizon/backend/driver/datastore/postgres"
	"github.com/rh-mithu/rizon/backend/internal/delivery/rest"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg    *config.Config
	server *http.Server
	db     *postgres.Store
}

func New(ctx context.Context, cfg *config.Config, l *slog.Logger) (*App, error) {
	handler := rest.ProvideHandler(cfg)
	server := &http.Server{
		Addr:    ":" + cfg.ServicePort,
		Handler: handler,
	}
	store := postgres.NewStore(cfg, l)
	return &App{
		cfg:    cfg,
		server: server,
		db:     store,
	}, nil
}

func (a *App) Run() error {
	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("Server starting on port %s", a.cfg.ServicePort)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("Start shutdown: %v", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := a.server.Shutdown(ctx); err != nil {
			err = a.server.Close()
			if err != nil {
				return err
			}
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func (a *App) Close() {
	if a.db != nil {
		err := a.db.DB().Close()
		if err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}
}
