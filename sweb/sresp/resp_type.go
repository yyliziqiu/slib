package sresp

import (
	"net/http"

	"github.com/yyliziqiu/slib/serror"
)

/**
错误码定义
A开头 客户端错误
B开头 服务端错误
C开头 三方服务错误

0-99      HTTP 协议定义的错误
100-999   框架定义的错误
1000-9999 用户自定义错误
*/

var (
	BadRequestError          = serror.New("A0001", "Bad Request").StatusCode(http.StatusBadRequest)
	UnauthorizedError        = serror.New("A0002", "Unauthorized").StatusCode(http.StatusUnauthorized)
	ForbiddenError           = serror.New("A0003", "Forbidden").StatusCode(http.StatusForbidden)
	NotFoundError            = serror.New("A0004", "Not Found").StatusCode(http.StatusNotFound)
	MethodNotAllowedError    = serror.New("A0005", "Method Not Allowed").StatusCode(http.StatusMethodNotAllowed)
	InternalServerErrorError = serror.New("B0001", "Internal Server Error").StatusCode(http.StatusInternalServerError)
)

type ErrorResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResult(code string, message string) ErrorResult {
	return ErrorResult{
		Code:    code,
		Message: message,
	}
}

func NewErrorResult2(err *serror.Error) ErrorResult {
	return NewErrorResult(err.Code, err.Message)
}

func errorResponse(err error, verbose bool) (int, ErrorResult) {
	var (
		statusCode = http.StatusBadRequest
		code       = BadRequestError.Code
		message    = BadRequestError.Message
	)

	zerr, ok := err.(*serror.Error)
	if ok {
		statusCode, code, message = zerr.Http()
	} else if verbose {
		message = err.Error()
	}

	return statusCode, NewErrorResult(code, message)
}
