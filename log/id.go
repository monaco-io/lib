package log

import (
	"context"

	"github.com/google/uuid"
)

type (
	contextKey    interface{}
	contextLogger struct {
		context.Context
	}
)

func WithContext(ctx *context.Context) *contextLogger {
	xRequestID := (*ctx).Value(contextKeyRequestID)
	if xRequestID == nil {
		*ctx = context.WithValue(*ctx, contextKeyRequestID, uuid.New().String())
	}
	return &contextLogger{*ctx}
}
