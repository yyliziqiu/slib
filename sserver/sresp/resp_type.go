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
	BadRequestError          = serror.New2(400, "A0001", "Bad Request")
	UnauthorizedError        = serror.New2(401, "A0002", "Unauthorized")
	ForbiddenError           = serror.New2(403, "A0003", "Forbidden")
	NotFoundError            = serror.New2(404, "A0004", "Not Found")
	MethodNotAllowedError    = serror.New2(405, "A0005", "Method Not Allowed")
	InternalServerErrorError = serror.New2(500, "B0001", "Internal Server Error")
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
