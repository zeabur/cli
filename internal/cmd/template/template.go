package template

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	templateCreateCmd "github.com/zeabur/cli/internal/cmd/template/create"
	templateDeleteCmd "github.com/zeabur/cli/internal/cmd/template/delete"
	templateDeployCmd "github.com/zeabur/cli/internal/cmd/template/deploy"
	templateGetCmd "github.com/zeabur/cli/internal/cmd/template/get"
	templateListCmd "github.com/zeabur/cli/internal/cmd/template/list"
	templateSearchCmd "github.com/zeabur/cli/internal/cmd/template/search"
	templateUpdateCmd "github.com/zeabur/cli/internal/cmd/template/update"
)

func NewCmdTemplate(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage templates",
	}

	cmd.AddCommand(templateListCmd.NewCmdList(f))
	cmd.AddCommand(templateDeployCmd.NewCmdDeploy(f))
	cmd.AddCommand(templateGetCmd.NewCmdGet(f))
	cmd.AddCommand(templateSearchCmd.NewCmdSearch(f))
	cmd.AddCommand(templateDeleteCmd.NewCmdDelete(f))
	cmd.AddCommand(templateCreateCmd.NewCmdCreate(f))
	cmd.AddCommand(templateUpdateCmd.NewCmdUpdate(f))

	return cmd
}
