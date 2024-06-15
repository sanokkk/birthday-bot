package logging

import (
	"log/slog"
	"os"
)

const (
	dev  = "dev"
	test = "test"
	prod = "prod"
)

func MustGetLogger(env string) slog.Logger {
	switch env {
	case dev:
		return *slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case test:
	case prod:
		return *slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	panic("unknown env")
}
