package main

import (
	"github.com/spf13/cobra"

	"changeme/internal/app"
)

// workerCmd 仅启动消息 Worker。
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "仅启动消息 Worker",
	RunE:  runWorker,
}

func init() {
	workerCmd.Flags().StringVar(&configFile, "config", "", "配置文件路径（默认 config.yaml）")
	rootCmd.AddCommand(workerCmd)
}

// runWorker 以 Worker 角色运行。
func runWorker(cmd *cobra.Command, _ []string) error {
	return app.Run(cmd.Context(), configFile, app.RoleWorker)
}
