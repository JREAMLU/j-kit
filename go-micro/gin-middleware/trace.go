package middleware

import (
	"github.com/JREAMLU/j-kit/go-micro/trace/opentracing"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
)

// HeaderTraceRespone set traceid in response header
func HeaderTraceRespone() gin.HandlerFunc {
	return func(c *gin.Context) {
		if md, ok := metadata.FromContext(c.Request.Context()); ok {
			c.Header("traceid", md[opentracing.ZipkinTraceID])
		}
		c.Next()
	}
}
