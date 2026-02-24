package list

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List dedicated servers",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f)
		},
	}

	return cmd
}

func runList(f *cmdutil.Factory) error {
	servers, err := f.ApiClient.GetServers(context.Background())
	if err != nil {
		return err
	}

	s := model.Servers(servers)
	f.Printer.Table(s.Header(), s.Rows())

	return nil
}
