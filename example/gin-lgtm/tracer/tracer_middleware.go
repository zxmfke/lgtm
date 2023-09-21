package tracer

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"net/http"
)

// Trace Gin的链路中间件
func Trace(ctx *gin.Context) {

	newSpan := Extract(ctx.Request, "trace-middleware")

	defer newSpan.End()

	ctx.Set(HeaderTraceContextKey, newSpan.spanCtx)
	ctx.Set(HeaderTraceIDKey, newSpan.SpanID())

	newSpan.Inject(ctx.Request)

	ctx.Next()

	newSpan.SetIntTag("http.status_code", ctx.Writer.Status())
	if ctx.Writer.Status() == http.StatusOK {
		newSpan.SetStatus(codes.Ok, "well done")
		return
	}
	newSpan.SetStatus(codes.Error, "something goes wrong")

}

func NewGinSpan(ctx *gin.Context, spanName string) Span {
	spanCtx, exist := ctx.Get(HeaderTraceContextKey)
	if !exist {
		return NewSpan(ctx.Request.Context(), spanName)
	}

	return NewSpan(spanCtx.(context.Context), spanName)
}
