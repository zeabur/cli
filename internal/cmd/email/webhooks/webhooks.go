package webhooks

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	webhooksCreateCmd "github.com/zeabur/cli/internal/cmd/email/webhooks/create"
	webhooksDeleteCmd "github.com/zeabur/cli/internal/cmd/email/webhooks/delete"
	webhooksGetCmd "github.com/zeabur/cli/internal/cmd/email/webhooks/get"
	webhooksListCmd "github.com/zeabur/cli/internal/cmd/email/webhooks/list"
	webhooksVerifyCmd "github.com/zeabur/cli/internal/cmd/email/webhooks/verify"
)

func NewCmdWebhooks(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "webhooks",
		Short:   "Manage email webhooks",
		Aliases: []string{"webhook"},
	}

	cmd.AddCommand(webhooksListCmd.NewCmdList(f))
	cmd.AddCommand(webhooksGetCmd.NewCmdGet(f))
	cmd.AddCommand(webhooksCreateCmd.NewCmdCreate(f))
	cmd.AddCommand(webhooksDeleteCmd.NewCmdDelete(f))
	cmd.AddCommand(webhooksVerifyCmd.NewCmdVerify(f))

	return cmd
}
