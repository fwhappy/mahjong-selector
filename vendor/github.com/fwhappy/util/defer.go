package util

import (
	"fmt"
	"runtime"
)

// RecoverPanic 捕获panic错误
func RecoverPanic() {
	if err := recover(); err != nil {
		timestamp := GetTimestamp()
		stack := make([]byte, 1024)
		stack = stack[:runtime.Stack(stack, true)]
		fmt.Println("[", timestamp, "]", "catchPanic:", err)
		fmt.Println("[", timestamp, "]", "stack:", string(stack))
	}
}
