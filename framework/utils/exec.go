package utils

import (
	"fmt"
	"os"
	"syscall"
)

// GetExecDirectory 获取当前执行程序目录
func GetExecDirectory() string {
	file, err := os.Getwd()
	if err == nil {
		return file + "/"
	}
	return ""
}
func CheckProcessExist(pid int) bool {
	//查询pid
	process, err := os.FindProcess(pid)
	if err != nil {

		return false
	}
	// 给进程发送signal 0, 如果返回nil，代表进程存在, 否则进程不存在
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
