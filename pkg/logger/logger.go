package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
}

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}
