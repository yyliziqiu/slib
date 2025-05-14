package sreq

import (
	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/slib/serror"
	"github.com/yyliziqiu/slib/sserver"
	"github.com/yyliziqiu/slib/sserver/sresp"
)

func bind(ctx *gin.Context, form interface{}, verbose bool) bool {
	err := ctx.ShouldBind(form)
	if err != nil {
		if logger := sserver.GetLogger(); logger != nil {
			logger.Warnf("Bind failed, path: %s, error: %v.", ctx.FullPath(), err)
		}
		if verbose {
			sresp.Error(ctx, serror.ParametersError.Wrap(err))
		} else {
			sresp.Error(ctx, serror.ParametersError)
		}
		return false
	}
	return true
}

func Bind(ctx *gin.Context, form interface{}) bool {
	return bind(ctx, form, false)
}

func BindVerbose(ctx *gin.Context, form interface{}) bool {
	return bind(ctx, form, true)
}
