package batch

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	batchGetCmd "github.com/zeabur/cli/internal/cmd/email/batch/get"
	batchListCmd "github.com/zeabur/cli/internal/cmd/email/batch/list"
)

func NewCmdBatch(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Manage batch email jobs",
	}

	cmd.AddCommand(batchListCmd.NewCmdList(f))
	cmd.AddCommand(batchGetCmd.NewCmdGet(f))

	return cmd
}
