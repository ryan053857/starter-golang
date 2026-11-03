package main

import "github.com/spf13/cobra"

// rootCmd 是应用根命令。
var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "golang-starter 应用",
}

// configFile 是 --config 指定的配置文件路径。
var configFile string
