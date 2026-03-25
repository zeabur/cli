package emails

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	emailsGetCmd "github.com/zeabur/cli/internal/cmd/email/emails/get"
	emailsListCmd "github.com/zeabur/cli/internal/cmd/email/emails/list"
)

func NewCmdEmails(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "emails",
		Short:   "View email sending records",
		Aliases: []string{"email-records", "records"},
	}

	cmd.AddCommand(emailsListCmd.NewCmdList(f))
	cmd.AddCommand(emailsGetCmd.NewCmdGet(f))

	return cmd
}
