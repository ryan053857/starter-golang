package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ryan053857/starter-golang/internal/app"
)

// httpCmd 仅启动 HTTP 服务。
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "仅启动 HTTP 服务",
	RunE:  runHTTP,
}

func init() {
	httpCmd.Flags().StringVar(&configFile, "config", "", "配置文件路径（默认 config.yaml）")
	rootCmd.AddCommand(httpCmd)
}

// runHTTP 以 HTTP 角色运行。
func runHTTP(cmd *cobra.Command, _ []string) error {
	if err := app.Run(cmd.Context(), configFile, app.RoleHTTP); err != nil {
		return fmt.Errorf("http: %w", err)
	}
	return nil
}
