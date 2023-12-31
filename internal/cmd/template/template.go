package template

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	templateDeleteCmd "github.com/zeabur/cli/internal/cmd/template/delete"
	templateDeployCmd "github.com/zeabur/cli/internal/cmd/template/deploy"
	templateGetCmd "github.com/zeabur/cli/internal/cmd/template/get"
	templateListCmd "github.com/zeabur/cli/internal/cmd/template/list"
)

func NewCmdTemplate(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage templates",
	}

	cmd.AddCommand(templateListCmd.NewCmdList(f))
	cmd.AddCommand(templateDeployCmd.NewCmdDeploy(f))
	cmd.AddCommand(templateGetCmd.NewCmdGet(f))
	cmd.AddCommand(templateDeleteCmd.NewCmdDelete(f))

	return cmd
}
