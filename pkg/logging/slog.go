package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/nzwice/surl/pkg/ctxutils"
)

const (
	timestampKey = "timestamp"
	ctxidKey     = "ctxid"
)

type hdl struct {
	next slog.Handler
}

// Enabled implements [slog.Handler].
func (h *hdl) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.next.Enabled(ctx, lvl)
}

// Handle implements [slog.Handler].
func (h *hdl) Handle(ctx context.Context, r slog.Record) error {
	if reqId := ctxutils.GetRequestId(ctx); reqId != "" {
		r.AddAttrs(slog.Any(ctxidKey, reqId))
	}
	return h.next.Handle(ctx, r)
}

// WithAttrs implements [slog.Handler].
func (h *hdl) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.next.WithAttrs(attrs)
}

// WithGroup implements [slog.Handler].
func (h *hdl) WithGroup(name string) slog.Handler {
	return h.next.WithGroup(name)
}

func SetupLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   timestampKey,
					Value: a.Value,
				}
			}
			return a
		},
	})
	decoratedHdl := &hdl{
		next: handler,
	}
	logger := slog.New(decoratedHdl)
	slog.SetDefault(logger)
}

func ErrorAttr(err any) slog.Attr {
	return slog.Any("error", err)
}
