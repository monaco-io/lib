package log

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_encoderConfig = zap.NewDevelopmentEncoderConfig()
	_autoLevel     = zap.NewAtomicLevel()

	name         string
	writer       zapcore.WriteSyncer
	errorWriter  zapcore.WriteSyncer
	caller       = zap.AddCaller()
	callerConfig = zap.AddCallerSkip(2)
	enc          = zapcore.NewJSONEncoder(_encoderConfig)
	trace        = zap.AddStacktrace(zap.ErrorLevel)

	core = zapcore.NewCore(zapcore.NewConsoleEncoder(_encoderConfig), os.Stdout, _autoLevel)
)

func RegisterWriter(writers ...io.Writer) {
	var writerTmp zapcore.WriteSyncer
	for _, w := range writers {
		if writerTmp == nil {
			writerTmp = zapcore.AddSync(w)
			continue
		}
		writerTmp = zap.CombineWriteSyncers(writerTmp, zapcore.AddSync(w))
	}
	writer = writerTmp
	newLogger()
}

func RegisterErrorWriter(writers ...io.Writer) {
	var writerTmp zapcore.WriteSyncer
	for _, w := range writers {
		if writerTmp == nil {
			writerTmp = zapcore.AddSync(w)
			continue
		}
		writerTmp = zap.CombineWriteSyncers(writerTmp, zapcore.AddSync(w))
	}
	errorWriter = writerTmp
	newLogger()
}

func RegisterDebug(debug bool) {
	if debug {
		_autoLevel.SetLevel(zap.DebugLevel)
		return
	}
	_autoLevel.SetLevel(zap.InfoLevel)
}

func RegisterServiceName(serviceName string) {
	name = serviceName
	newLogger()
}
