package types

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

// Result .
type Result struct {
	Code int                   `json:"code"` // 0 means ok, else means error
	Meta map[string]HealthInfo `json:"meta"`
}

type alias struct {
	checker HealthChecker
	name    string
	timeout time.Duration
}

// CheckingMgr .
type CheckingMgr struct {
	checkers []alias
}

// NewHealthMgr .
func NewHealthMgr() *CheckingMgr {
	return &CheckingMgr{
		checkers: make([]alias, 0),
	}
}

const defaultTimeout = 10 * time.Second

// AddChecker of mgr
func (mgr *CheckingMgr) AddChecker(name string, hchecker HealthChecker, timeout int) {
	a := alias{
		checker: hchecker,
		name:    name,
		timeout: time.Duration(timeout) * time.Second,
	}
	if timeout <= 0 {
		a.timeout = defaultTimeout
	}
	mgr.checkers = append(mgr.checkers, a)
}

func (mgr *CheckingMgr) doCheck() Result {
	wg := sync.WaitGroup{}
	wg.Add(len(mgr.checkers))
	ch := make(chan HealthInfo, len(mgr.checkers))

	for _, a := range mgr.checkers {
		ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
		go func(ctx context.Context, a alias) {
			var info = NewHealthInfo()

			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					info.name = a.name
					info.Healthy = false
					info.Meta["panic"] = r
					info.Meta["stack"] = string(debug.Stack())
					ch <- info
				}
			}()

			select {
			case <-ctx.Done():
				// append timeout error
				info.name = a.name
				info.Healthy = false
				info.Meta["error"] = "check timeout"
			default:
				info = a.checker.Check()
				info.name = a.name
			}

			ch <- info
			cancel()
		}(ctx, a)
	}

	wg.Wait()
	close(ch)
	// merge check info
	var result Result
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
func (mgr *CheckingMgr) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result := mgr.doCheck()
		byts, err := json.Marshal(result)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Fprintf(w, string(byts))
	}
}

// GinHandler .
func (mgr *CheckingMgr) GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		result := mgr.doCheck()
		c.JSON(http.StatusOK, result)
	}
}
