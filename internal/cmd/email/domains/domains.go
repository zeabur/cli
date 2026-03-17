package domains

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	domainsAddCmd "github.com/zeabur/cli/internal/cmd/email/domains/add"
	domainsDeleteCmd "github.com/zeabur/cli/internal/cmd/email/domains/delete"
	domainsGetCmd "github.com/zeabur/cli/internal/cmd/email/domains/get"
	domainsListCmd "github.com/zeabur/cli/internal/cmd/email/domains/list"
	domainsVerifyCmd "github.com/zeabur/cli/internal/cmd/email/domains/verify"
)

func NewCmdDomains(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domains",
		Short:   "Manage email domains",
		Aliases: []string{"domain"},
	}

	cmd.AddCommand(domainsListCmd.NewCmdList(f))
	cmd.AddCommand(domainsAddCmd.NewCmdAdd(f))
	cmd.AddCommand(domainsGetCmd.NewCmdGet(f))
	cmd.AddCommand(domainsVerifyCmd.NewCmdVerify(f))
	cmd.AddCommand(domainsDeleteCmd.NewCmdDelete(f))

	return cmd
}
