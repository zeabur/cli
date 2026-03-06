package provider

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdProvider(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "providers",
		Short:   "List available dedicated server providers",
		Args:    cobra.NoArgs,
		Aliases: []string{"provider"},
		RunE: func(cmd *cobra.Command, args []string) error {
			providers, err := f.ApiClient.ListDedicatedServerProviders(context.Background())
			if err != nil {
				return fmt.Errorf("list providers failed: %w", err)
			}

			if len(providers) == 0 {
				f.Log.Infof("No providers available")
				return nil
			}

			header := []string{"Code", "Name"}
			rows := make([][]string, len(providers))
			for i, p := range providers {
				rows[i] = []string{p.Code, p.Name}
			}
			f.Printer.Table(header, rows)

			return nil
		},
	}

	return cmd
}
