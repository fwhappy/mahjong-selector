package mselector

import (
	"fmt"
)

// 显示错误
func showError(format string, s ...interface{}) {
	fmt.Println("[ERROR]", fmt.Sprintf(format, s...))
}

// 显示debug信息
func showDebug(format string, s ...interface{}) {
	fmt.Println("[DEBUG]", fmt.Sprintf(format, s...))
}
