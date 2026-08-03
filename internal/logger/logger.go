package logger

import (
	"io"
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init(config Config) {
	var out io.Writer = os.Stdout
	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.Level == slog.LevelDebug}

	var handler slog.Handler

	if config.Json {
		handler = slog.NewJSONHandler(out, opts)

	} else {
		handler = slog.NewTextHandler(out, opts)
	}

	Log = slog.New(handler)
}
