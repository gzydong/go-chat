package errorx

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int
	Message string
}

func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func IsError(e error) bool {
	return errors.As(e, new(*Error))
}
