package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	debug = iota
	info
	warn
	err
	panic
	fatal
	dev
)

const arrow = "->"

var (
	contextKeyRequestID contextKey = "x-request-id"
)

var (
	log *zap.Logger // core
)

func init() {
	RegisterWriter()
}

func newLogger() {
	if writer != nil {
		writerCore := zapcore.NewCore(enc, writer, _autoLevel)
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
