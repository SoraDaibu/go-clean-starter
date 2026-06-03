// Package logger provides the application's structured logger as an interface
// backed by the standard library's log/slog.
//
// Exposing a Logger interface (instead of a concrete *slog.Logger) lets layers
// depend on the abstraction rather than the backend (DIP); the slog handler is
// swappable without touching call sites. With(args ...any) Logger lets each
// caller derive its own scoped logger with additional fields.
package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger is the structured-logging contract used across the application.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	// With returns a Logger that includes the given key/value pairs on every
	// subsequent log line, so callers can scope their own fields.
	With(args ...any) Logger
}

// slogLogger is the slog-backed implementation of Logger.
type slogLogger struct {
	l *slog.Logger
}

func (s *slogLogger) Debug(msg string, args ...any) { s.l.Debug(msg, args...) }
func (s *slogLogger) Info(msg string, args ...any)  { s.l.Info(msg, args...) }
func (s *slogLogger) Warn(msg string, args ...any)  { s.l.Warn(msg, args...) }
func (s *slogLogger) Error(msg string, args ...any) { s.l.Error(msg, args...) }

func (s *slogLogger) With(args ...any) Logger {
	return &slogLogger{l: s.l.With(args...)}
}

// New builds a Logger for the given log level and environment.
//
// The "local" env uses a human-readable text handler to ease development; any
// other env emits JSON so logs are machine-parsable in production. Source
// location (file:line) is included to match the previous logger's behavior.
func New(level, env string) Logger {
	opts := &slog.HandlerOptions{
		Level:     parseLevel(level),
		AddSource: true,
	}

	var handler slog.Handler
	if env == "local" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return &slogLogger{l: slog.New(handler)}
}

// parseLevel maps a configured level string to a slog.Level, defaulting to Info.
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warning", "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// SetDefault makes l the default logger for slog's package-level functions, so
// bootstrap/CLI code that logs via slog shares the same configured backend.
func SetDefault(l Logger) {
	if sl, ok := l.(*slogLogger); ok {
		slog.SetDefault(sl.l)
	}
}

// ctxKey is an unexported context key for storing a request-scoped logger.
type ctxKey struct{}

// IntoContext returns a copy of ctx carrying the given logger.
func IntoContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext returns the logger stored in ctx, falling back to a logger that
// wraps slog.Default() when none is present.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return l
	}
	return &slogLogger{l: slog.Default()}
}
