package template

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	templateDeployCmd "github.com/zeabur/cli/internal/cmd/template/deploy"
	templateListCmd "github.com/zeabur/cli/internal/cmd/template/list"
)

func NewCmdTemplate(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage templates",
	}

	cmd.AddCommand(templateListCmd.NewCmdList(f))
	cmd.AddCommand(templateDeployCmd.NewCmdDeploy(f))

	return cmd
}
