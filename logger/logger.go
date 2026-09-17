package logger

import (
	"connect-sample/requestid"
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
)

var _ slog.Handler = (*Handler)(nil)

type Handler struct {
	handler slog.Handler
}

func NewHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	return &Handler{
		handler: slog.NewJSONHandler(w, opts),
	}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	// デバッグ: コンテキストに入っている値と型を確認
	raw := ctx.Value(counterKey)
	if raw == nil {
		fmt.Println("debug/ Handler.Handle: ctx counterKey = nil")
	} else {
		fmt.Printf("debug/ Handler.Handle: ctx counterKey type=%T value=%v\n", raw, raw)
	}

	if pc, ok := raw.(uintptr); ok && pc != 0 {
		fmt.Printf("debug/ Handler.Handle: override r.PC=%#x fn=%s\n", pc, runtime.FuncForPC(pc).Name())
		r.PC = pc
	} else {
		fmt.Printf("debug/ Handler.Handle: use r.PC=%#x fn=%s\n", r.PC, runtime.FuncForPC(r.PC).Name())
	}

	return h.handler.Handle(ctx, r)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		handler: h.handler.WithGroup(name),
	}
}

type counterKeyType struct{}

var counterKey = counterKeyType{}

func WithPC(ctx context.Context, skipCount int) context.Context {
	var pcs [1]uintptr
	runtime.Callers(skipCount, pcs[:])
	return context.WithValue(ctx, counterKey, pcs[0])
}

const skippingFrameCount = 3

func Info(ctx context.Context, msg string, args ...any) {
	l := requestid.FromContext(ctx)
	if l == nil {
		slog.Default().Log(WithPC(ctx, skippingFrameCount), slog.LevelInfo, msg, args...)
		return
	}
	l.Log(WithPC(ctx, skippingFrameCount), slog.LevelInfo, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	requestid.FromContext(WithPC(ctx, skippingFrameCount)).Error(msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	requestid.FromContext(WithPC(ctx, skippingFrameCount)).Debug(msg, args...)
}
