package xxml

import (
	"bytes"
	"encoding/xml"
	"io"
)

func MarshalStringX(v any) string {
	b, _ := xml.Marshal(v)
	return string(b)
}

func MarshalString(v any) (string, error) {
	b, err := xml.Marshal(v)
	return string(b), err
}

func Unmarshal(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

func UnmarshalString(str string, v any) error {
	return xml.Unmarshal([]byte(str), v)
}

func UnmarshalT[T any](b []byte) (T, error) {
	var v T
	err := xml.Unmarshal(b, &v)
	return v, err
}

func UnmarshalStringT[T any](str string) (T, error) {
	var v T
	err := xml.Unmarshal([]byte(str), &v)
	return v, err
}

func Marshal[T any](data T) ([]byte, error) {
	bytes, err := xml.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func MarshalX[T any](data T) []byte {
	bytes, err := xml.Marshal(data)
	if err != nil {
		return nil
	}
	return bytes
}

func MarshalReaderX[T any](data T) io.Reader {
	return bytes.NewReader(MarshalX(data))
}

func MarshalReader[T any](data T) (io.Reader, error) {
	xmlBytes, err := Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(xmlBytes), nil
}
