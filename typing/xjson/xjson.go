package xjson

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/monaco-io/lib/typing/xopt"
	"github.com/monaco-io/lib/typing/xstr"
)

func marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func Marshal(v any) ([]byte, error) {
	return marshal(v)
}

func MarshalX(v any) []byte {
	b, _ := marshal(v)
	return b
}

func MarshalIndent(v any, opts ...xopt.Option[marshalIndentOption]) ([]byte, error) {
	var opt marshalIndentOption
	for _, o := range opts {
		o(&opt)
	}
	return json.MarshalIndent(v, opt.prefix, opt.indent)
}

type marshalIndentOption struct {
	prefix string
	indent string
}

func WithIndentPrefix(prefix string) xopt.Option[marshalIndentOption] {
	return func(o *marshalIndentOption) {
		o.prefix = prefix
	}
}
func WithIndentString(indent string) xopt.Option[marshalIndentOption] {
	return func(o *marshalIndentOption) {
		o.indent = indent
	}
}

func MarshalIndentX(v any, opts ...xopt.Option[marshalIndentOption]) []byte {
	var opt marshalIndentOption
	for _, o := range opts {
		o(&opt)
	}
	b, _ := json.MarshalIndent(v, opt.prefix, opt.indent)
	return b
}

func MarshalIndentStringX(v any, opts ...xopt.Option[marshalIndentOption]) string {
	var opt marshalIndentOption
	opt.indent = xstr.TAB
	for _, o := range opts {
		o(&opt)
	}
	b, _ := json.MarshalIndent(v, opt.prefix, opt.indent)
	return string(b)
}

func MarshalStringX(v any) string {
	b, _ := marshal(v)
	return string(b)
}

func MarshalString(v any) (string, error) {
	b, err := marshal(v)
	return string(b), err
}

func Unmarshal(data []byte, v any) error {
	return unmarshal(data, v)
}

func UnmarshalString(str string, v any) error {
	return unmarshal([]byte(str), v)
}

func UnmarshalT[T any](b []byte) (T, error) {
	var v T
	err := unmarshal(b, &v)
	return v, err
}

func UnmarshalStringT[T any](str string) (T, error) {
	var v T
	err := unmarshal([]byte(str), &v)
	return v, err
}

func MarshalReader[T any](data T) (io.Reader, error) {
	jsonBytes, err := Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonBytes), nil
}

func MarshalReaderX[T any](data T) io.Reader {
	return bytes.NewReader(MarshalX(data))
}

func TransformT[Target, From any](from From) (*Target, error) {
	b, err := marshal(from)
	if err != nil {
		return nil, err
	}
	return UnmarshalT[*Target](b)
}
