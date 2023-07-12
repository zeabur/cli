package context

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	contextGetCmd "github.com/zeabur/cli/internal/cmd/context/get"
	contextSetCmd "github.com/zeabur/cli/internal/cmd/context/set"
)

func NewCmdContext(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "context <command>",
		Short:   "Manage contexts",
		Aliases: []string{"ctx"},
	}

	cmd.AddCommand(contextGetCmd.NewCmdGet(f))
	cmd.AddCommand(contextSetCmd.NewCmdSet(f))
	cmd.AddCommand(contextGetCmd.NewCmdGet(f))

	return cmd
}
