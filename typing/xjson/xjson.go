package xjson

import (
	"bytes"
	"encoding/json"
	"io"
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

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func MarshalIndentX(v any, prefix, indent string) []byte {
	b, _ := json.MarshalIndent(v, prefix, indent)
	return b
}

func MarshalIndentStringX(v any, prefix, indent string) string {
	b, _ := json.MarshalIndent(v, prefix, indent)
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
