package autorenew

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id      string
	enable  bool
	disable bool
}

func NewCmdAutoRenew(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "auto-renew",
		Short: "Toggle auto-renew for a registered domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAutoRenew(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Registered domain ID")
	cmd.Flags().BoolVar(&opts.enable, "enable", false, "Enable auto-renew")
	cmd.Flags().BoolVar(&opts.disable, "disable", false, "Disable auto-renew")

	return cmd
}

func runAutoRenew(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()

	if opts.enable && opts.disable {
		return fmt.Errorf("cannot use both --enable and --disable")
	}

	if opts.id == "" {
		if !f.Interactive {
			return fmt.Errorf("--id is required")
		}
		domains, err := f.ApiClient.ListRegisteredDomains(ctx)
		if err != nil {
			return fmt.Errorf("list registered domains failed: %w", err)
		}
		if len(domains) == 0 {
			return fmt.Errorf("no registered domains found")
		}

		options := make([]string, len(domains))
		for i, d := range domains {
			autoRenew := "off"
			if d.AutoRenew {
				autoRenew = "on"
			}
			options[i] = fmt.Sprintf("%s (auto-renew: %s)", d.Domain, autoRenew)
		}
		idx, err := f.Prompter.Select("Select domain", "", options)
		if err != nil {
			return err
		}
		opts.id = domains[idx].ID
	}

	if !opts.enable && !opts.disable {
		if !f.Interactive {
			return fmt.Errorf("--enable or --disable is required")
		}
		idx, err := f.Prompter.Select("Auto-renew setting", "", []string{"Enable", "Disable"})
		if err != nil {
			return err
		}
		if idx == 0 {
			opts.enable = true
		} else {
			opts.disable = true
		}
	}

	autoRenew := opts.enable

	domain, err := f.ApiClient.SetDomainAutoRenew(ctx, opts.id, autoRenew)
	if err != nil {
		return fmt.Errorf("set auto-renew failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(domain)
	}

	status := "disabled"
	if domain.AutoRenew {
		status = "enabled"
	}
	f.Log.Infof("Auto-renew %s for %s", status, domain.Domain)
	return nil
}
