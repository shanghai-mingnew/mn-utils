package logger

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GinLoggerMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		TimeStamp := time.Now()
		Latency := TimeStamp.Sub(start)

		ClientIP := c.ClientIP()
		Method := c.Request.Method
		StatusCode := c.Writer.Status()
		ErrorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		Path := path

		if len(ErrorMessage) > 0 {
			Errorf("[GIN] %3d | %13v | %15s |%-7s %#v\n%s",
				StatusCode,
				Latency,
				ClientIP,
				Method,
				Path,
				ErrorMessage)
		} else if StatusCode >= http.StatusBadRequest {
			Errorf("[GIN] %3d | %13v | %15s |%-7s %#v",
				StatusCode,
				Latency,
				ClientIP,
				Method,
				Path)
		} else {
			Debugf("[GIN] %3d | %13v | %15s |%-7s %#v",
				StatusCode,
				Latency,
				ClientIP,
				Method,
				Path)
		}
	}
}
