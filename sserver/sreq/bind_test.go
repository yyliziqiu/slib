package sreq

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type TestForm struct {
	Name string `form:"name" binding:"required"`
}

func TestBind(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	var form TestForm
	ok := Bind(ctx, &form)
	if !ok {
		t.Error("bind failed")
	}
}
