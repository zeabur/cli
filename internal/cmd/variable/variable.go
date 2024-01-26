package variable

import (
	"github.com/spf13/cobra"

	variableListCmd "github.com/zeabur/cli/internal/cmd/variable/list"
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

	return cmd
}
