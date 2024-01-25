package domain

import (
	"github.com/spf13/cobra"

	domainCreateCmd "github.com/zeabur/cli/internal/cmd/domain/create"
	domainDeleteCmd "github.com/zeabur/cli/internal/cmd/domain/delete"
	domainListCmd "github.com/zeabur/cli/internal/cmd/domain/list"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdDomain(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain",
		Short:   "Manage domains",
		Long:    `Manage domains`,
		Aliases: []string{"domain"},
	}

	cmd.AddCommand(domainCreateCmd.NewCmdCreateDomain(f))
	cmd.AddCommand(domainListCmd.NewCmdListDomains(f))
	cmd.AddCommand(domainDeleteCmd.NewCmdDeleteDomain(f))

	return cmd
}
