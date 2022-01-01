package xec

import (
	"sync/atomic"
)

var (
	_messages atomic.Value // NOTE: stored map[int]map[Language]string
	_lang     Language     = LangDefault
)

type Message map[int]map[Language]string

// Register register ecode message map.
func Register(code int, msg map[Language]string) {
	if msg == nil {
		return
	}
	var mp Message
	if m, ok := _messages.Load().(Message); !ok || m == nil {
		mp = make(Message)
	}
	if mp[code] == nil {
		mp[code] = make(map[Language]string)
	}

	mp[code] = msg
	_messages.Store(mp)
}
