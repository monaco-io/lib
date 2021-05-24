package log

func _handler(c *contextLogger, level Level, msg string, keyValues ...interface{}) {
	sugar := log.Sugar()
	if c != nil {
		sugar = sugar.With(contextKeyRequestID.(string), (*c.Context).Value(contextKeyRequestID).(string))
	}
	defer func() { _ = log.Sync() }()

	switch level {
	case _debug:
		sugar.Debugw(msg, keyValues...)
	case _info:
		sugar.Infow(msg, keyValues...)
	case _warn:
		sugar.Warnw(msg, keyValues...)
	case _err:
		sugar.Errorw(msg, keyValues...)
	case _panic:
		sugar.Panicw(msg, keyValues...)
	case _fatal:
		sugar.Fatalw(msg, keyValues...)
	case _dev:
		sugar.Infow(msg, keyValues...)
	}
}
