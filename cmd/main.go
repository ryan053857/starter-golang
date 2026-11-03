// Package main 是应用入口与 cobra 子命令。
package main

import "os"

// main 执行根命令；失败时以状态码 1 退出。
func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
