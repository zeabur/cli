package dns

import (
	"github.com/spf13/cobra"

	dnsCreateCmd "github.com/zeabur/cli/internal/cmd/domain/dns/create"
	dnsDeleteCmd "github.com/zeabur/cli/internal/cmd/domain/dns/delete"
	dnsListCmd "github.com/zeabur/cli/internal/cmd/domain/dns/list"
	dnsUpdateCmd "github.com/zeabur/cli/internal/cmd/domain/dns/update"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdDNS(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns <command>",
		Short: "Manage DNS records for registered domains",
	}

	cmd.AddCommand(dnsListCmd.NewCmdList(f))
	cmd.AddCommand(dnsCreateCmd.NewCmdCreate(f))
	cmd.AddCommand(dnsUpdateCmd.NewCmdUpdate(f))
	cmd.AddCommand(dnsDeleteCmd.NewCmdDelete(f))

	return cmd
}
