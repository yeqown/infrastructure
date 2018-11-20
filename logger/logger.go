package logger

import (
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

// Logger wrap logrus.Logger
type Logger struct {
	*logrus.Logger
}

// Println to implement interface `gorm.LogWriter`
func (L Logger) Println(format string, v ...interface{}) {
	L.Logger.Infof(format, v...)
}

var (
	// Log is default logger only output to os.Stdout
	Log = logrus.New()
)

// initialize default logger
func init() {
	Log.Formatter = &logrus.JSONFormatter{}
}

// NewJSONLogger ...
// init logger to set filepath and format
func NewJSONLogger(logPath, filename, lv string) (*Logger, error) {
	logger := logrus.New()
	logger.AddHook(NewHook())

	logger.Formatter = &logrus.JSONFormatter{}

	level, err := logrus.ParseLevel(lv)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)

	fd, err := os.OpenFile(
		path.Join(logPath, filename),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		if os.IsNotExist(err) {
			if os.MkdirAll(logPath, 0777) != nil {
				return nil, err
			}
			goto Finally
		}
		return nil, err
	}

Finally:
	logger.Out = io.MultiWriter(os.Stdout, fd)
	return &Logger{logger}, nil
}
