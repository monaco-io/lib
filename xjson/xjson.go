package xjson

import "github.com/bytedance/sonic"

var xjson = sonic.Config{
	NoQuoteTextMarshaler:    true,
	NoValidateJSONMarshaler: true,
	NoValidateJSONSkip:      true,
}.Froze()

func Marshal(v any) ([]byte, error) {
	return xjson.Marshal(v)
}

func MarshalX(v any) []byte {
	b, _ := xjson.Marshal(v)
	return b
}

func MarshalStringX(v any) string {
	b, _ := xjson.Marshal(v)
	return string(b)
}

func MarshalString(v any) (string, error) {
	return xjson.MarshalToString(v)
}

func Unmarshal(data []byte, v any) error {
	return xjson.Unmarshal(data, v)
}

func UnmarshalString(str string, v any) error {
	return xjson.UnmarshalFromString(str, v)
}

func UnmarshalT[T any](b []byte) (T, error) {
	var v T
	err := xjson.Unmarshal(b, &v)
	return v, err
}

func UnmarshalStringT[T any](str string) (T, error) {
	var v T
	err := xjson.UnmarshalFromString(str, &v)
	return v, err
}
