package flags

import (
	"io"
	"log/slog"
	"os"

	"github.com/titusjaka/go-sample/internal/infrastructure/nopslog"
)

type Logger struct {
	Level  string `kong:"optional,group='Logger',name=log-level,env=LOG_LEVEL,enum='debug,info,warn,error,off',default=info,help='The minimal level for the logs (${enum}).'"`
	Format string `kong:"optional,group='Logger',name=log-format,env=LOG_FORMAT,enum='text,json',default='json',help='The format of the log output (${enum}).'"`
}

func (f Logger) Init() *slog.Logger {
	if f.Level == "off" {
		return nopslog.NewNoplogger()
	}

	opts := &slog.HandlerOptions{
		Level: f.SlogLevel(),
	}

	handler := f.SlogHandler(os.Stdout, opts)

	return slog.New(handler)
}

func (f Logger) SlogLevel() slog.Level {
	switch f.Level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (f Logger) SlogHandler(out io.Writer, opts *slog.HandlerOptions) slog.Handler {
	switch f.Format {
	case "json":
		return slog.NewJSONHandler(out, opts)
	case "text":
		return slog.NewTextHandler(out, opts)
	default:
		return slog.NewJSONHandler(out, opts)
	}
}
