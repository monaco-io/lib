package xec

type Error interface {
	Message() string
	Code() int
}
