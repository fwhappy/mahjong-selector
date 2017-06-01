package util

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
)

// RandIntn 获取一个 0 ~ n 之间的随机值
func RandIntn(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Intn(n)
}

// GetRandString 生成n位随机数字字符串
func GetRandString(n int) string {
	var buffer bytes.Buffer
	for i := 0; i < n; i++ {
		buffer.WriteString(strconv.Itoa(RandIntn(10)))
	}

	return buffer.String()
}
