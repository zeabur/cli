package email

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	emailBatchCmd "github.com/zeabur/cli/internal/cmd/email/batch"
	emailDomainsCmd "github.com/zeabur/cli/internal/cmd/email/domains"
	emailEmailsCmd "github.com/zeabur/cli/internal/cmd/email/emails"
	emailKeysCmd "github.com/zeabur/cli/internal/cmd/email/keys"
	emailScheduledCmd "github.com/zeabur/cli/internal/cmd/email/scheduled"
	emailSendCmd "github.com/zeabur/cli/internal/cmd/email/send"
	emailStatusCmd "github.com/zeabur/cli/internal/cmd/email/status"
	emailWebhooksCmd "github.com/zeabur/cli/internal/cmd/email/webhooks"
)

func NewCmdEmail(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "email",
		Short: "Manage Zeabur Email",
	}

	cmd.AddCommand(emailStatusCmd.NewCmdStatus(f))
	cmd.AddCommand(emailDomainsCmd.NewCmdDomains(f))
	cmd.AddCommand(emailKeysCmd.NewCmdKeys(f))
	cmd.AddCommand(emailWebhooksCmd.NewCmdWebhooks(f))
	cmd.AddCommand(emailSendCmd.NewCmdSend(f))
	cmd.AddCommand(emailEmailsCmd.NewCmdEmails(f))
	cmd.AddCommand(emailScheduledCmd.NewCmdScheduled(f))
	cmd.AddCommand(emailBatchCmd.NewCmdBatch(f))

	return cmd
}
