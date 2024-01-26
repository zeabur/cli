package variable

import (
	"github.com/spf13/cobra"

	variableCreateCmd "github.com/zeabur/cli/internal/cmd/variable/create"
	varableDeleteCmd "github.com/zeabur/cli/internal/cmd/variable/delete"
	variableListCmd "github.com/zeabur/cli/internal/cmd/variable/list"
	variableUpdateCmd "github.com/zeabur/cli/internal/cmd/variable/update"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdVariable(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "variable",
		Short:   "Manage environment variables",
		Long:    `Manage environment variables of service`,
		Aliases: []string{"var"},
	}

	cmd.AddCommand(variableListCmd.NewCmdListVariables(f))
	cmd.AddCommand(variableCreateCmd.NewCmdCreateVariable(f))
	cmd.AddCommand(variableUpdateCmd.NewCmdUpdateVariable(f))
	cmd.AddCommand(varableDeleteCmd.NewCmdDeleteVariable(f))

	return cmd
}
