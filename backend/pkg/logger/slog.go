package logger

import (
	"github.com/rh-mithu/rizon/backend/config"
	"log/slog"
	"os"
	"time"
)

func NewSlog(cfg *config.Config) *slog.Logger {
	isProduction := cfg.Env == "production"

	var handler slog.Handler

	opts := &slog.HandlerOptions{
		// This includes the file and line number in the log output
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(a.Value.Time().Format(time.RFC3339)),
				}
			}
			return a
		},
	}

	if isProduction {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Overwrite level for development
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
