package log

func D(msg string, keyValues ...interface{}) {
	_handler(nil, debug, msg, keyValues...)
}

func I(msg string, keyValues ...interface{}) {
	_handler(nil, info, msg, keyValues...)
}

func W(msg string, keyValues ...interface{}) {
	_handler(nil, warn, msg, keyValues...)
}

func E(msg string, keyValues ...interface{}) {
	_handler(nil, err, msg, keyValues...)
}

func P(msg string, keyValues ...interface{}) {
	_handler(nil, panic, msg, keyValues...)
}

func F(msg string, keyValues ...interface{}) {
	_handler(nil, fatal, msg, keyValues...)
}

func Log(keyValues ...interface{}) {
	_handler(nil, dev, arrow, keyValues...)
}

func (c *contextLogger) D(msg string, keyValues ...interface{}) {
	_handler(c, debug, msg, keyValues...)
}

func (c *contextLogger) I(msg string, keyValues ...interface{}) {
	_handler(c, info, msg, keyValues...)
}

func (c *contextLogger) W(msg string, keyValues ...interface{}) {
	_handler(c, warn, msg, keyValues...)
}

func (c *contextLogger) E(msg string, keyValues ...interface{}) {
	_handler(c, err, msg, keyValues...)
}

func (c *contextLogger) P(msg string, keyValues ...interface{}) {
	_handler(c, panic, msg, keyValues...)
}

func (c *contextLogger) F(msg string, keyValues ...interface{}) {
	_handler(c, fatal, msg, keyValues...)
}

func (c *contextLogger) Log(keyValues ...interface{}) {
	_handler(c, dev, arrow, keyValues...)
}
