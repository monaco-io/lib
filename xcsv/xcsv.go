package xcsv

import (
	"os"

	"github.com/gocarina/gocsv"
)

func Unmarshal[T any](data []byte, out []T) error {
	return gocsv.UnmarshalBytes(data, out)
}

func UnmarshalT[T any](data []byte) (out []T, err error) {
	err = gocsv.UnmarshalBytes(data, &out)
	return
}

func MarshalBytes[T any](data []T) ([]byte, error) {
	return gocsv.MarshalBytes(data)
}

func MarshalBytesX[T any](data []T) []byte {
	bytes, _ := gocsv.MarshalBytes(data)
	return bytes
}

func WriteFile[T any](data []T, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	return gocsv.MarshalFile(data, file)
}
