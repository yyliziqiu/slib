package sresp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ============ Response ============

func Response(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, data)
}

func ResponseError(ctx *gin.Context, statusCode int, code string, message string) {
	ctx.JSON(statusCode, NewErrorResult(code, message))
}

func OK(ctx *gin.Context) {
	ctx.String(http.StatusOK, "")
}

func Result(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func Error(ctx *gin.Context, err error) {
	ctx.JSON(errorResponse(err, false))
}

func ErrorVerbose(ctx *gin.Context, err error) {
	ctx.JSON(errorResponse(err, true))
}

func ErrorString(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, NewErrorResult(BadRequestError.Code, message))
}

// ============ Abort ============

func AbortOK(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

func AbortResult(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, data)
}

func AbortError(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(errorResponse(err, false))
}

func AbortErrorVerbose(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(errorResponse(err, true))
}

func AbortErrorString(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResult(BadRequestError.Code, message))
}

// ============ Handle ============

func AbortBadRequest(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(errorResponse(BadRequestError, false))
}

func AbortUnauthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(errorResponse(UnauthorizedError, false))
}

func AbortForbidden(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(errorResponse(ForbiddenError, false))
}

func AbortNotFound(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(errorResponse(NotFoundError, false))
}

func AbortMethodNotAllowed(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(errorResponse(MethodNotAllowedError, false))
}

func AbortInternalServerError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(errorResponse(InternalServerErrorError, false))
}
