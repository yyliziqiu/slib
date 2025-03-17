package sreq

import (
	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/slib/serror"
	"github.com/yyliziqiu/slib/sserver"
	"github.com/yyliziqiu/slib/sserver/sresp"
)

var (
	ParametersError = serror.New2(400, "A0100", "parameters error")
)

func bind(ctx *gin.Context, form interface{}, verbose bool) bool {
	err := ctx.ShouldBind(form)
	if err != nil {
		logger := sserver.GetLogger()
		if logger != nil {
			logger.Errorf("Bind failed, path: %s, error: %v.", ctx.FullPath(), err)
		}
		if verbose {
			sresp.Error(ctx, ParametersError.Wrap(err))
		} else {
			sresp.Error(ctx, ParametersError)
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
