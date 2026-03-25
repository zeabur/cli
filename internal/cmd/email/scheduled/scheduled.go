package scheduled

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	scheduledCancelCmd "github.com/zeabur/cli/internal/cmd/email/scheduled/cancel"
	scheduledGetCmd "github.com/zeabur/cli/internal/cmd/email/scheduled/get"
	scheduledListCmd "github.com/zeabur/cli/internal/cmd/email/scheduled/list"
)

func NewCmdScheduled(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduled",
		Short: "Manage scheduled emails",
	}

	cmd.AddCommand(scheduledListCmd.NewCmdList(f))
	cmd.AddCommand(scheduledGetCmd.NewCmdGet(f))
	cmd.AddCommand(scheduledCancelCmd.NewCmdCancel(f))

	return cmd
}
