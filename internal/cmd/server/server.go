package server

import (
	"github.com/spf13/cobra"

	serverCatalogCmd "github.com/zeabur/cli/internal/cmd/server/catalog"
	serverGetCmd "github.com/zeabur/cli/internal/cmd/server/get"
	serverListCmd "github.com/zeabur/cli/internal/cmd/server/list"
	serverPlanCmd "github.com/zeabur/cli/internal/cmd/server/plan"
	serverProviderCmd "github.com/zeabur/cli/internal/cmd/server/provider"
	serverRebootCmd "github.com/zeabur/cli/internal/cmd/server/reboot"
	serverRegionCmd "github.com/zeabur/cli/internal/cmd/server/region"
	serverRentCmd "github.com/zeabur/cli/internal/cmd/server/rent"
	serverSSHCmd "github.com/zeabur/cli/internal/cmd/server/ssh"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdServer(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server <command>",
		Short: "Manage dedicated servers",
	}

	cmd.AddCommand(serverCatalogCmd.NewCmdCatalog(f))
	cmd.AddCommand(serverListCmd.NewCmdList(f))
	cmd.AddCommand(serverGetCmd.NewCmdGet(f))
	cmd.AddCommand(serverRebootCmd.NewCmdReboot(f))
	cmd.AddCommand(serverRentCmd.NewCmdRent(f))
	cmd.AddCommand(serverProviderCmd.NewCmdProvider(f))
	cmd.AddCommand(serverRegionCmd.NewCmdRegion(f))
	cmd.AddCommand(serverPlanCmd.NewCmdPlan(f))
	cmd.AddCommand(serverSSHCmd.NewCmdSSH(f))

	return cmd
}
