package typing

import (
	"github.com/monaco-io/lib/typing/xec"
	"github.com/monaco-io/lib/typing/xjson"
)

type ResponseEntity struct {
	xec.Error
	Data any `json:"data" xml:"data" yaml:"data"`
}

func (r *ResponseEntity) ToString() string {
	return xjson.MarshalStringX(r)
}
