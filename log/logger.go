package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int

const (
	_debug Level = iota
	_info
	_warn
	_err
	_panic
	_fatal
	_dev
)

var levelNames = [...]string{
	_debug: "DEBUG",
	_info:  "INFO",
	_warn:  "WARN",
	_err:   "ERROR",
	_panic: "PANIC",
	_fatal: "FATAL",
	_dev:   "DEV",
}

// String implementation.
func (l Level) String() string {
	return levelNames[l]
}

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
