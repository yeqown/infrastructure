package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Hook ... define logrus hook type
type Hook struct {
	Field     string
	Skip      int
	levels    []logrus.Level
	Formatter func(file, function string, line int) string
}

// Levels ...
func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

// Fire ...
func (hook *Hook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = hook.Formatter(findCaller(hook.Skip))
	return nil
}

// NewHook ...
func NewHook(levels ...logrus.Level) *Hook {
	hook := &Hook{
		Field:  "file",
		Skip:   5,
		levels: levels,
		Formatter: func(file, function string, line int) string {
			return fmt.Sprintf("%s:%d", file, line)
		},
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}

	return hook
}

func findCaller(skip int) (string, string, int) {
	var (
		pc       uintptr
		file     string
		function string
		line     int
	)
	for i := 0; i < 10; i++ {
		pc, file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	if pc != 0 {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		function = frame.Function
	}

	return file, function, line
}

func getCaller(skip int) (uintptr, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return 0, "", 0
	}

	splited := strings.Split(file, "/")
	file = strings.Join(splited[len(splited)-2:], "/")

	return pc, file, line
}
