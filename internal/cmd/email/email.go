package email

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	emailDomainsCmd "github.com/zeabur/cli/internal/cmd/email/domains"
	emailKeysCmd "github.com/zeabur/cli/internal/cmd/email/keys"
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

	return cmd
}
