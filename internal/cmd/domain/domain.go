package domain

import (
	"github.com/spf13/cobra"
	
	domainCreateCmd "github.com/zeabur/cli/internal/cmd/domain/create"
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

	return cmd
}
