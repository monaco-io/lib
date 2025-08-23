package xlog

import (
	"context"

	"github.com/monaco-io/lib/typing/xstr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level zapcore.Level

var log *zap.Logger // core

func init() {
	RegisterWriter()
}

func newLogger() {
	if writer != nil {
		writerCore := zapcore.NewCore(enc, writer, _level)
		core = zapcore.NewTee(core, writerCore)
		writer = nil
	}
	if errorWriter != nil {
		writerCore := zapcore.NewCore(enc, errorWriter, zap.ErrorLevel)
		core = zapcore.NewTee(core, writerCore)
		errorWriter = nil
	}
	log = zap.New(core, caller, callerConfig, trace).Named(name)
}

func _handler(ctx context.Context, level zapcore.Level, msg string, keyValues ...interface{}) {
	sugar := log.Sugar()
	if ctx != nil {
		id, ok := ctx.Value(xstr.X_REQUEST_ID).(string)
		if ok {
			sugar = sugar.With(xstr.X_REQUEST_ID, id)
		}
	}
	sugar.Logw(level, msg, keyValues...)
}
