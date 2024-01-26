package variable

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdVariable(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "variable",
		Short:   "Manage environment variables",
		Long:    `Manage environment variables of service`,
		Aliases: []string{"var"},
	}

	return cmd
}
