package main

import (
	"github.com/spf13/cobra"

	"github.com/runforyou-ai/golang-starter/internal/app"
)

// serveCmd 启动全部角色（HTTP、Worker 等）。
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 HTTP、Worker 等全部角色",
	RunE:  runServe,
}

func init() {
	serveCmd.Flags().StringVar(&configFile, "config", "", "配置文件路径（默认 config.yaml）")
	rootCmd.AddCommand(serveCmd)
}

// runServe 以全部角色运行。
func runServe(cmd *cobra.Command, _ []string) error {
	return app.Run(cmd.Context(), configFile, app.RoleAll)
}
