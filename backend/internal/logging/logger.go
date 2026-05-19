package logging

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	level := slog.LevelInfo
	if os.Getenv("KANVIX_LOG_LEVEL") != "" {
		if err := level.UnmarshalText([]byte(os.Getenv("KANVIX_LOG_LEVEL"))); err != nil {
			level = slog.LevelInfo
		}
	}
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(h)
}
