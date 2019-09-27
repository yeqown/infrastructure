package healthcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthInfo to describe the status of target
type HealthInfo struct {
	name    string
	Healthy bool                   `json:"healthy"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// NewHealthInfo .
func NewHealthInfo() HealthInfo {
	return HealthInfo{
		Healthy: false,
		Meta:    make(map[string]interface{}),
	}
}

// HealthChecker .
type HealthChecker interface {
	Check() HealthInfo
}

// httpResponse .
type httpResponse struct {
	// Code 0 means ok, else means error
	Code int `json:"code"`
	// Meta includes all checkers' `HealthInfo` result
	Meta map[string]HealthInfo `json:"meta"`
}

type checkerWrapper struct {
	// origin checker to do the actual checking work.
	checker HealthChecker
	// name to the wrapper to represent the output field too.
	name string
	// timeout duation of the checking work is supposed to be.
	timeout time.Duration
}

// HealthCheckingMgr .
type HealthCheckingMgr struct {
	checkerWrappers []checkerWrapper
}

// NewHealthMgr .
func NewHealthMgr() *HealthCheckingMgr {
	return &HealthCheckingMgr{
		checkerWrappers: make([]checkerWrapper, 0),
	}
}

const defaultTimeout = 10 * time.Second

// AddChecker of mgr
func (mgr *HealthCheckingMgr) AddChecker(name string, hchecker HealthChecker, timeout int) {
	a := checkerWrapper{
		checker: hchecker,
		name:    name,
		timeout: time.Duration(timeout) * time.Second,
	}
	if timeout <= 0 {
		a.timeout = defaultTimeout
	}
	mgr.checkerWrappers = append(mgr.checkerWrappers, a)
}

func (mgr *HealthCheckingMgr) doCheck(ctx context.Context) httpResponse {
	wg := sync.WaitGroup{}
	wg.Add(len(mgr.checkerWrappers))
	ch := make(chan HealthInfo, len(mgr.checkerWrappers))

	for _, c := range mgr.checkerWrappers {
		tctx, cancel := context.WithTimeout(ctx, c.timeout)
		go func(ctx context.Context, w checkerWrapper) {
			var info = NewHealthInfo()
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					info.name = w.name
					info.Healthy = false
					info.Meta["panic"] = r
					info.Meta["stack"] = string(debug.Stack())
					ch <- info
				}
			}()

			select {
			case <-ctx.Done():
				// append timeout error
				info.name = w.name
				info.Healthy = false
				info.Meta["error"] = "check timeout"
			default:
				info = w.checker.Check()
				info.name = w.name
			}

			ch <- info
			cancel()
		}(tctx, c)
	}

	wg.Wait()
	close(ch)
	// merge check info
	var result httpResponse
	result.Meta = make(map[string]HealthInfo)
	for v := range ch {
		if !v.Healthy && result.Code == 0 {
			result.Code = -1
		}
		result.Meta[v.name] = v
	}
	return result
}

// Handler .
func (mgr *HealthCheckingMgr) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result := mgr.doCheck(req.Context())
		byts, err := json.Marshal(result)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Fprintf(w, string(byts))
	}
}

// GinHandler .
func (mgr *HealthCheckingMgr) GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		result := mgr.doCheck(c.Request.Context())
		c.JSON(http.StatusOK, result)
	}
}
