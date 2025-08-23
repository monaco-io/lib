package xstring

import (
	"strings"

	"github.com/google/uuid"
)

func Pick(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func UUID() string {
	return uuid.NewString()
}

func UUIDX() string {
	return strings.ReplaceAll(UUID(), "-", "")
}
