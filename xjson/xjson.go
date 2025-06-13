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

func MarshalString(v any) (string, error) {
	return xjson.MarshalToString(v)
}

func Unmarshal(data []byte, v any) error {
	return xjson.Unmarshal(data, v)
}

func UnmarshalString(str string, v any) error {
	return xjson.UnmarshalFromString(str, v)
}
