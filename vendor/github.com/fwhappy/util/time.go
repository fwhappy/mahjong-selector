package util

import (
	"fmt"
	"time"
)

// GetTime 获取当前时间戳
func GetTime() int64 {
	return time.Now().Unix()
}

// GetMicrotime 获取微秒时间
func GetMicrotime() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GetTimestamp 获取当前格式化时间
func GetTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// FormatUnixTime 将时间戳格式化
func FormatUnixTime(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("2006-01-02 15:04:05")
}

// GetYMD 获取当前年月日
func GetYMD() string {
	return time.Now().Format("20060102")
}

// GetChinaWeekDay 获取中国的星期几
func GetChinaWeekDay() int {
	weekDay := time.Now().Weekday()
	if weekDay == time.Sunday {
		return 7
	}
	return int(weekDay)
}
