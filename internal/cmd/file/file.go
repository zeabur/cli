package file

import (
	"github.com/spf13/cobra"

	fileListCmd "github.com/zeabur/cli/internal/cmd/file/list"
	fileReadCmd "github.com/zeabur/cli/internal/cmd/file/read"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdFile(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file <command>",
		Short: "Manage uploaded files",
	}

	cmd.AddCommand(fileListCmd.NewCmdList(f))
	cmd.AddCommand(fileReadCmd.NewCmdRead(f))

	return cmd
}
