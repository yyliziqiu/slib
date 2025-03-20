package sresp

import (
	"net/http"

	"github.com/yyliziqiu/slib/serror"
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
		status  = http.StatusBadRequest
		code    = serror.BadRequest.Code
		message = serror.BadRequest.Message
	)

	zerr, ok := err.(*serror.Error)
	if ok {
		status, code, message = zerr.Http()
	} else if verbose {
		message = err.Error()
	}

	return status, NewErrorResult(code, message)
}
