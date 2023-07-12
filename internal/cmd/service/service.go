package service

import (
	"github.com/spf13/cobra"

	serviceGetCmd "github.com/zeabur/cli/internal/cmd/service/get"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdService(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service <command>",
		Short:   "Manage services",
		Aliases: []string{"svc"},
	}

	cmd.AddCommand(serviceGetCmd.NewCmdGet(f))

	return cmd
}
