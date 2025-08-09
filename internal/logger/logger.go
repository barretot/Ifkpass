package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init(env string) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	Log = slog.New(handler)

	Log.Info("logger initialized", "env", env)
}
