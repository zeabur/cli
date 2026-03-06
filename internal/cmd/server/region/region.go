package region

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	provider string
}

func NewCmdRegion(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "regions",
		Short:   "List available regions for a dedicated server provider",
		Args:    cobra.NoArgs,
		Aliases: []string{"region"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRegion(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.provider, "provider", "", "Provider code (required)")

	return cmd
}

func runRegion(f *cmdutil.Factory, opts *Options) error {
	if opts.provider == "" {
		if f.Interactive {
			providers, err := f.ApiClient.ListDedicatedServerProviders(context.Background())
			if err != nil {
				return fmt.Errorf("list providers failed: %w", err)
			}
			if len(providers) == 0 {
				return fmt.Errorf("no providers available")
			}
			options := make([]string, len(providers))
			for i, p := range providers {
				options[i] = p.Name
			}
			idx, err := f.Prompter.Select("Select a provider", "", options)
			if err != nil {
				return err
			}
			opts.provider = providers[idx].Code
		} else {
			return fmt.Errorf("--provider is required")
		}
	}

	regions, err := f.ApiClient.ListDedicatedServerRegions(context.Background(), opts.provider)
	if err != nil {
		return fmt.Errorf("list regions failed: %w", err)
	}

	if len(regions) == 0 {
		f.Log.Infof("No regions available for provider %s", opts.provider)
		return nil
	}

	header := []string{"ID", "Name", "City", "Country", "Continent"}
	rows := make([][]string, len(regions))
	for i, r := range regions {
		rows[i] = []string{r.ID, r.Name, r.City, r.Country, r.Continent}
	}
	f.Printer.Table(header, rows)

	return nil
}
