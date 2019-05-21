package ginic

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	logger "github.com/yeqown/infrastructure/framework/logrus-logger"
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
		rbw := &respBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = rbw

		start := time.Now()
		ctxCpy := c.Copy()

		c.Next()

		latency := time.Now().Sub(start)
		fields := make(map[string]interface{})
		fields["requestData"] = parseRequestForm(ctxCpy)
		if logResponse {
			fields["responseData"] = rbw.body.String()
		}

		Logger.WithFields(fields).Infof("[Request] %v |%3d| %13v | %15s |%-7s %s",
			start.Format("2006/01/02 - 15:04:05"),
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}

func parseRequestForm(ctxCpy *gin.Context) (form map[string]interface{}) {
	form = make(map[string]interface{})

	switch ctxCpy.Request.Method {
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

// CORS ...
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
