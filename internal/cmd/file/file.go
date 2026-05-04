package file

import (
	"github.com/spf13/cobra"

	filePullCmd "github.com/zeabur/cli/internal/cmd/file/pull"
	"github.com/zeabur/cli/internal/cmdutil"
)

// NewCmdFile creates the file command.
func NewCmdFile(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file <command>",
		Short: "Manage uploaded files",
	}

	cmd.AddCommand(filePullCmd.NewCmdPull(f))

	return cmd
}
