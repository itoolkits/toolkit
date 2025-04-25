// init logger

package ctxt

import (
	"io"
	"log/slog"
)

var AccessLogger = slog.Default()

// InitSystemLogger - init system logger
func InitSystemLogger(w io.Writer) {
	logHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}

// InitAccessLogger - init access logger
func InitAccessLogger(w io.Writer) {
	AccessLogger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}))
}
