package sserver

import (
	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/slib/serror"
	"github.com/yyliziqiu/slib/sserver/sresp"
)

var ParamError = serror.New("A0100", "request params error")

func BindForm(ctx *gin.Context, form interface{}, verbose bool) bool {
	err := ctx.ShouldBind(form)
	if err != nil {
		if _errorLogger != nil {
			_errorLogger.Warnf("Bind request params failed, path: %s, form: %v, error: %v.", ctx.FullPath(), form, err)
		}
		if verbose {
			sresp.Error(ctx, ParamError.Wrap(err))
		} else {
			sresp.Error(ctx, ParamError)
		}
		return false
	}
	return true
}
