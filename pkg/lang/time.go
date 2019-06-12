package lang

import (
	"time"
)

const (
	// Layout1 = "2006-01-02 15:04:06"
	Layout1 = "2006-01-02 15:04:06"
	// Layout2 = "20060102"
	Layout2 = "20060102"
)

// CurTimeFormat get current time string format value
func CurTimeFormat(layout string) string {
	return time.Now().Format(layout)
}

// ParseTime format time string into a time.Time with specified layout
func ParseTime(layout, value string) time.Time {
	t, _ := time.Parse(layout, value)
	return t
}

// GetTodayDate 获取今天的日期
func GetTodayDate() time.Time {
	now := time.Now()
	return GetTimeDate(now)
}

// GetTimeDate 获取给定时间的日期整数
func GetTimeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}
