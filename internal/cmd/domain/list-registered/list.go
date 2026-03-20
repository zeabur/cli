package listregistered

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdListRegistered(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-registered",
		Short:   "List purchased domains",
		Args:    cobra.NoArgs,
		Aliases: []string{"lr"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f)
		},
	}
	return cmd
}

func runList(f *cmdutil.Factory) error {
	domains, err := f.ApiClient.ListRegisteredDomains(context.Background())
	if err != nil {
		return fmt.Errorf("list registered domains failed: %w", err)
	}

	if len(domains) == 0 {
		if f.JSON {
			return f.Printer.JSON([]any{})
		}
		f.Log.Infof("No registered domains found")
		return nil
	}

	if f.JSON {
		return f.Printer.JSON(domains)
	}

	f.Printer.Table(domains.Header(), domains.Rows())
	return nil
}
