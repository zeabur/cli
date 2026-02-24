// Package server contains the cmd for managing dedicated servers
package server

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	serverListCmd "github.com/zeabur/cli/internal/cmd/server/list"
)

// NewCmdServer creates the server command
func NewCmdServer(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage dedicated servers",
	}

	cmd.AddCommand(serverListCmd.NewCmdList(f))

	return cmd
}
