package serror

import (
	"fmt"
	"net/http"
)

type Error struct {
	statusCode int

	Code    string
	Message string
}

func New(code string, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func (e *Error) clone(message string) *Error {
	return &Error{
		statusCode: e.statusCode,
		Code:       e.Code,
		Message:    message,
	}
}

func (e *Error) With(v interface{}) *Error {
	message := ""
	switch v.(type) {
	case error:
		message = v.(error).Error()
	case string:
		message = v.(string)
	default:
		message = fmt.Sprintf("%v", v)
	}
	return e.clone(message)
}

func (e *Error) Wrap(err error) *Error {
	return e.clone(fmt.Sprintf("%s [%v]", e.Message, err))
}

func (e *Error) Format(message string, a ...interface{}) *Error {
	return e.clone(fmt.Sprintf(message, a...))
}

func (e *Error) Fields(a ...interface{}) *Error {
	return e.clone(fmt.Sprintf(e.Message, a...))
}

func (e *Error) StatusCode(code int) *Error {
	e.statusCode = code
	return e
}

func (e *Error) Http() (int, string, string) {
	if e.statusCode != 0 {
		return e.statusCode, e.Code, e.Message
	}

	statusCode := http.StatusBadRequest
	if e.Code[0] != 'A' {
		statusCode = http.StatusInternalServerError
	}

	return statusCode, e.Code, e.Message
}
