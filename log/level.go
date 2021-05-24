package log

func D(msg string, keyValues ...interface{}) {
	_handler(nil, _debug, msg, keyValues...)
}

func I(msg string, keyValues ...interface{}) {
	_handler(nil, _info, msg, keyValues...)
}

func W(msg string, keyValues ...interface{}) {
	_handler(nil, _warn, msg, keyValues...)
}

func E(msg string, keyValues ...interface{}) {
	_handler(nil, _err, msg, keyValues...)
}

func P(msg string, keyValues ...interface{}) {
	_handler(nil, _panic, msg, keyValues...)
}

func F(msg string, keyValues ...interface{}) {
	_handler(nil, _fatal, msg, keyValues...)
}

func Log(keyValues ...interface{}) {
	_handler(nil, _dev, arrow, keyValues...)
}

func (c *contextLogger) D(msg string, keyValues ...interface{}) {
	_handler(c, _debug, msg, keyValues...)
}

func (c *contextLogger) I(msg string, keyValues ...interface{}) {
	_handler(c, _info, msg, keyValues...)
}

func (c *contextLogger) W(msg string, keyValues ...interface{}) {
	_handler(c, _warn, msg, keyValues...)
}

func (c *contextLogger) E(msg string, keyValues ...interface{}) {
	_handler(c, _err, msg, keyValues...)
}

func (c *contextLogger) P(msg string, keyValues ...interface{}) {
	_handler(c, _panic, msg, keyValues...)
}

func (c *contextLogger) F(msg string, keyValues ...interface{}) {
	_handler(c, _fatal, msg, keyValues...)
}

func (c *contextLogger) Log(keyValues ...interface{}) {
	_handler(c, _dev, arrow, keyValues...)
}
