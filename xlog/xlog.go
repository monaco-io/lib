package xlog

import (
	"context"

	"go.uber.org/zap/zapcore"
)

func Sync() error {
	if log == nil {
		return nil
	}
	return log.Sync()
}

func D(ctx context.Context, msg string, keyValues ...interface{}) {
	_handler(ctx, zapcore.DebugLevel, msg, keyValues...)
}

func I(ctx context.Context, msg string, keyValues ...interface{}) {
	_handler(ctx, zapcore.InfoLevel, msg, keyValues...)
}

func W(ctx context.Context, msg string, keyValues ...interface{}) {
	_handler(ctx, zapcore.WarnLevel, msg, keyValues...)
}

func E(ctx context.Context, msg string, keyValues ...interface{}) {
	_handler(ctx, zapcore.ErrorLevel, msg, keyValues...)
}

func P(ctx context.Context, msg string, keyValues ...interface{}) {
	_handler(ctx, zapcore.PanicLevel, msg, keyValues...)
}

func F(ctx context.Context, msg string, keyValues ...interface{}) {
	_handler(ctx, zapcore.FatalLevel, msg, keyValues...)
}
