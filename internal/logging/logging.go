package logging

import (
	"context"
	"io"
	"log"
	"log/slog"
)

// Logger returns an opinionated logger for the given level.
func Logger(out io.Writer, level slog.Level) *slog.Logger {
	return slog.New(textHandler(out, level))
}

func textHandler(out io.Writer, level slog.Level) slog.Handler {
	return slog.NewTextHandler(out, &slog.HandlerOptions{
		Level: level,
	})
}

// TODO: need this?
func textLogger(out io.Writer) *log.Logger {
	return log.New(out, "", 0)
}

// DiscardLogger returns a [slog.Logger] that discards all output.
func DiscardLogger() *slog.Logger {
	return slog.New(DiscardHandler())
}

// DiscardHandler returns a [slog.Handler] that discards all output.
// It is an implementation of https://github.com/golang/go/issues/62005.
func DiscardHandler() slog.Handler {
	return (*discardHandler)(nil)
}

type discardHandler struct {
	slog.Handler
}

func (*discardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
