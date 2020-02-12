package ginic

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	logger "github.com/yeqown/infrastructure/framework/logrus-logger"
	"github.com/yeqown/infrastructure/pkg/session"
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

const (
	// SessionKey .
	SessionKey = "session"
	// HTTPTokenHeaderName .
	HTTPTokenHeaderName = "X-Token"

	bRefreshTokenAutomatic = false // 是否主动更新token过期时间
)

// FactoryToIToken .
type FactoryToIToken func() session.IToken

// AuthJWT .
func AuthJWT(mgr session.ITokenManager, factory FactoryToIToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokString := c.Request.Header.Get(HTTPTokenHeaderName)
		if tokString == "" {
			// true: 没有在请求头中携带token头，提示错误
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": -4,
				"msg":  "用户认证标识错误",
			})
			return
		}

		headerTok := factory()
		if err := mgr.PraseToken(tokString, headerTok); err != nil {
			// true: 解析失败或者token内容非法
			logger.Log.WithField("token", tokString).Warnf("mgr.PraseToken() failed, err=%v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": -4,
				"msg":  "用户认证标识非法",
			})
			return
		}

		// 是否过期
		cacheTok := factory()
		if err := mgr.Get(headerTok.TokenKey(), cacheTok); err != nil {
			if err == redis.Nil {
				// true: key missed
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code": -5,
					"msg":  "用户认证标识已过期，请重新登录",
				})
				return
			}

			// redis 异常
			logger.Log.WithField("token", tokString).Warnf("mgr.PraseToken() failed to cmp with redis, err=%v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
			return
		}

		if bRefreshTokenAutomatic {
			// true: 主动更新token过期时间
			// 刷新token过期时间（可选）
			if cacheTok.Expiration() < session.OneWeekExpired {
				cacheTok.SetExpiration(session.OneMonthExpired)
				if err := mgr.Refresh(cacheTok, session.OneMonthExpired); err != nil {
					logger.Log.Warnf("mgr.Refresh() failed, err=%v", err)
				}
			}
		}

		c.Set(SessionKey, cacheTok)
	}
}
