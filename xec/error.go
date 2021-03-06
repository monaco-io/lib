package xec

type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string {
	return e.Message
}

func New(code int, msg string) Error {
	return Error{
		Code:    code,
		Message: msg,
	}
}
