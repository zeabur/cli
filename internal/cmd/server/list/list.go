package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List dedicated servers",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f)
		},
	}

	return cmd
}

func runList(f *cmdutil.Factory) error {
	servers, err := f.ApiClient.ListServers(context.Background())
	if err != nil {
		return fmt.Errorf("list servers failed: %w", err)
	}

	if len(servers) == 0 {
		f.Log.Infof("No servers found")
		return nil
	}

	f.Printer.Table(servers.Header(), servers.Rows())

	return nil
}
