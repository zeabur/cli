// Package service provides the entry point of the service command
package service

import (
	"github.com/spf13/cobra"

	serviceDeployCmd "github.com/zeabur/cli/internal/cmd/service/deploy"
	serviceExposeCmd "github.com/zeabur/cli/internal/cmd/service/expose"
	serviceGetCmd "github.com/zeabur/cli/internal/cmd/service/get"
	serviceListCmd "github.com/zeabur/cli/internal/cmd/service/list"
	serviceMetricCmd "github.com/zeabur/cli/internal/cmd/service/metric"
	serviceRedeployCmd "github.com/zeabur/cli/internal/cmd/service/redeploy"
	serviceRestartCmd "github.com/zeabur/cli/internal/cmd/service/restart"
	serviceSuspendCmd "github.com/zeabur/cli/internal/cmd/service/suspend"
	"github.com/zeabur/cli/internal/cmdutil"
)

// NewCmdService creates the service command
func NewCmdService(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service <command>",
		Short:   "Manage services",
		Aliases: []string{"svc"},
	}

	cmd.AddCommand(serviceGetCmd.NewCmdGet(f))
	cmd.AddCommand(serviceListCmd.NewCmdList(f))
	cmd.AddCommand(serviceExposeCmd.NewCmdExpose(f))
	cmd.AddCommand(serviceMetricCmd.NewCmdMetric(f))
	cmd.AddCommand(serviceRestartCmd.NewCmdRestart(f))
	cmd.AddCommand(serviceRedeployCmd.NewCmdRedeploy(f))
	cmd.AddCommand(serviceSuspendCmd.NewCmdSuspend(f))
	cmd.AddCommand(serviceDeployCmd.NewCmdDeploy(f))

	return cmd
}
