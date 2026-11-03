package main

import (
	"github.com/spf13/cobra"

	"github.com/runforyou-ai/golang-starter/internal/app"
)

// serveCmd 启动 HTTP 服务。
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 HTTP 服务",
	RunE:  runServe,
}

func init() {
	serveCmd.Flags().StringVar(&configFile, "config", "", "配置文件路径（默认 config.yaml）")
	rootCmd.AddCommand(serveCmd)
}

// runServe 启动 HTTP 服务直至退出。
func runServe(cmd *cobra.Command, _ []string) error {
	return app.Run(cmd.Context(), configFile)
}
