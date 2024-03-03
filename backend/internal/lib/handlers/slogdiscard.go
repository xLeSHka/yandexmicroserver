package slogdiscard

import (
	"context"
	"log/slog"
)

/*игнорирует все сообщения которые мы в него отправляем*/
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	/*игнорируем запись журнала*/
	return nil
}
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
