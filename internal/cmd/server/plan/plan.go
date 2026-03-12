package plan

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	provider string
	region   string
}

func NewCmdPlan(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "plans",
		Short:   "List available plans for a dedicated server provider and region",
		Args:    cobra.NoArgs,
		Aliases: []string{"plan"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPlan(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.provider, "provider", "", "Provider code (required)")
	cmd.Flags().StringVar(&opts.region, "region", "", "Region ID (required)")

	return cmd
}

func runPlan(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()

	if opts.provider == "" {
		if f.Interactive {
			providers, err := f.ApiClient.ListDedicatedServerProviders(ctx)
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

	if opts.region == "" {
		if f.Interactive {
			regions, err := f.ApiClient.ListDedicatedServerRegions(ctx, opts.provider)
			if err != nil {
				return fmt.Errorf("list regions failed: %w", err)
			}
			if len(regions) == 0 {
				return fmt.Errorf("no regions available for provider %s", opts.provider)
			}
			options := make([]string, len(regions))
			for i, r := range regions {
				options[i] = fmt.Sprintf("%s (%s, %s)", r.Name, r.City, r.Country)
			}
			idx, err := f.Prompter.Select("Select a region", "", options)
			if err != nil {
				return err
			}
			opts.region = regions[idx].ID
		} else {
			return fmt.Errorf("--region is required")
		}
	}

	plans, err := f.ApiClient.ListDedicatedServerPlans(ctx, opts.provider, opts.region)
	if err != nil {
		return fmt.Errorf("list plans failed: %w", err)
	}

	if len(plans) == 0 {
		if f.JSON {
			return f.Printer.JSON([]any{})
		}
		f.Log.Infof("No plans available for provider %s in region %s", opts.provider, opts.region)
		return nil
	}

	if f.JSON {
		return f.Printer.JSON(plans)
	}

	f.Printer.Table(plans.Header(), plans.Rows())

	return nil
}
