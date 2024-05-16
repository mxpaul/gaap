package loggy

import (
	"log/slog"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var ErrLogLevelInvalid = errors.New("log level string value invalid")

type Config struct {
	Level string `yaml:"level,omitempty"`
}

func NewLogger(cfg Config) (*slog.Logger, error) {
	level, err := slogLevelFromString(cfg.Level)
	if err != nil {
		return nil, errors.Wrapf(ErrLogLevelInvalid, "got value: %s", cfg.Level)
	}

	logOpt := slog.HandlerOptions{Level: level}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &logOpt))
	return logger, nil
}

func slogLevelFromString(s string) (l slog.Level, err error) {
	switch strings.ToLower(s) {
	case "debug":
		l = slog.LevelDebug
	case "":
		fallthrough
	case "info":
		l = slog.LevelInfo
	case "warning":
		fallthrough
	case "warn":
		l = slog.LevelWarn
	case "err":
		fallthrough
	case "error":
		l = slog.LevelError
	default:
		return l, ErrLogLevelInvalid
	}
	return
}
