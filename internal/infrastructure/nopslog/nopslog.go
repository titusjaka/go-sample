package nopslog

import (
	"context"
	"log/slog"
)

type Handler struct{}

var _ slog.Handler = Handler{}

func (Handler) Enabled(context.Context, slog.Level) bool { return false }

func (Handler) Handle(context.Context, slog.Record) error { return nil }

func (h Handler) WithAttrs([]slog.Attr) slog.Handler { return h }

func (h Handler) WithGroup(string) slog.Handler { return h }

func NewNoplogger() *slog.Logger {
	return slog.New(Handler{})
}
