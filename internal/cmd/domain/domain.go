package domain

import (
	"github.com/spf13/cobra"

	autoRenewCmd "github.com/zeabur/cli/internal/cmd/domain/auto-renew"
	domainCreateCmd "github.com/zeabur/cli/internal/cmd/domain/create"
	domainDeleteCmd "github.com/zeabur/cli/internal/cmd/domain/delete"
	dnsCmd "github.com/zeabur/cli/internal/cmd/domain/dns"
	getRegisteredCmd "github.com/zeabur/cli/internal/cmd/domain/get-registered"
	domainListCmd "github.com/zeabur/cli/internal/cmd/domain/list"
	listRegisteredCmd "github.com/zeabur/cli/internal/cmd/domain/list-registered"
	purchaseCmd "github.com/zeabur/cli/internal/cmd/domain/purchase"
	registrantCmd "github.com/zeabur/cli/internal/cmd/domain/registrant"
	renewCmd "github.com/zeabur/cli/internal/cmd/domain/renew"
	searchCmd "github.com/zeabur/cli/internal/cmd/domain/search"
	verificationCmd "github.com/zeabur/cli/internal/cmd/domain/verification"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdDomain(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain",
		Short:   "Manage domains",
		Long:    `Manage domains`,
		Aliases: []string{"domain"},
	}

	// Existing service domain binding commands
	cmd.AddCommand(domainCreateCmd.NewCmdCreateDomain(f))
	cmd.AddCommand(domainListCmd.NewCmdListDomains(f))
	cmd.AddCommand(domainDeleteCmd.NewCmdDeleteDomain(f))

	// Domain registration commands
	cmd.AddCommand(searchCmd.NewCmdSearch(f))
	cmd.AddCommand(purchaseCmd.NewCmdPurchase(f))
	cmd.AddCommand(listRegisteredCmd.NewCmdListRegistered(f))
	cmd.AddCommand(getRegisteredCmd.NewCmdGetRegistered(f))
	cmd.AddCommand(renewCmd.NewCmdRenew(f))
	cmd.AddCommand(autoRenewCmd.NewCmdAutoRenew(f))
	cmd.AddCommand(dnsCmd.NewCmdDNS(f))
	cmd.AddCommand(registrantCmd.NewCmdRegistrant(f))
	cmd.AddCommand(verificationCmd.NewCmdVerification(f))

	return cmd
}
