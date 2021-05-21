package log

func _handler(c *contextLogger, level int, msg string, keyValues ...interface{}) {
	sugar := log.Sugar()
	if c != nil {
		sugar = sugar.With(contextKeyRequestID.(string), c.Context.Value(contextKeyRequestID).(string))
	}
	defer func() { _ = log.Sync() }()

	switch level {
	case debug:
		sugar.Debugw(msg, keyValues...)
	case info:
		sugar.Infow(msg, keyValues...)
	case warn:
		sugar.Warnw(msg, keyValues...)
	case err:
		sugar.Errorw(msg, keyValues...)
	case panic:
		sugar.Panicw(msg, keyValues...)
	case fatal:
		sugar.Fatalw(msg, keyValues...)
	case dev:
		sugar.Infow(msg, keyValues...)
	}
}
