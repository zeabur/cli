package deployment

import (
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"

	deploymentGetCmd "github.com/zeabur/cli/internal/cmd/deployment/get"
	deploymentListCmd "github.com/zeabur/cli/internal/cmd/deployment/list"
)

func NewCmdDeployment(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deployment",
		Short:   "Manage deployments",
		Long:    `Manage deployments`,
		Aliases: []string{"deploy"},
	}

	cmd.AddCommand(deploymentListCmd.NewCmdList(f))
	cmd.AddCommand(deploymentGetCmd.NewCmdGet(f))

	return cmd
}
