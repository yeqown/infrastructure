package ginic

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yeqown/server-common/logger"
)

// Recovery is a middleware to record each panic into file
// usage like:
//	gin.Engine.Use(Recovery(*os.File))
func Recovery(out io.Writer) gin.HandlerFunc {
	// self custom RecoveryWithWriter rather than gin.RecoveryWithWriter(io.Writer)
	return RecoveryWithWriter(out)
}

type respBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w respBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LogRequest is a middleware to log each request
func LogRequest(Logger *logger.Logger, logResponse bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer := &respBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		start := time.Now()
		ctxCpy := c.Copy()

		c.Next()

		latency := time.Now().Sub(start)
		fields := make(map[string]interface{})
		fields["requestData"] =  parseRequestForm(ctxCpy)
		if logResponse {
			fields["responseData"] = rbw.body.String(),
		}

		Logger.WithFields(fields).Infof("[Request] %v |%3d| %13v | %15s |%-7s %s",
			end.Format("2006/01/02 - 15:04:05"),
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}

func parseRequestForm(c *gin.Context) (form map[string]interface{}) {
	switch c.Request.Method {
	case http.MethodPost, http.MethodPut:
		ctxCpy.Request.ParseMultipartForm(32 << 20)
	case http.MethodGet:
		ctxCpy.Request.ParseForm()
	default:
		ctxCpy.Request.ParseForm()
	}

	for k, v := range ctxCpy.Request.Form {
		form[k] = v
	}

	return
}