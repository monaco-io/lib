package xec

import (
	"errors"
	"fmt"
)

func New(code int, msg ...string) Error {
	return Error{
		Code:    code,
		Message: fmt.Sprintf("%v", msg),
	}
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`

	cause error
}

func (e Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e Error) Wrap(err error) Error {
	return Error{
		Code:    e.Code,
		Message: e.Message,
		cause:   err,
	}
}

func (e Error) Is(err error) bool {
	if err == nil {
		return false
	}
	var e2 Error
	if errors.As(err, &e2) {
		return e.Code == e2.Code
	}
	return e.Error() == err.Error()
}
