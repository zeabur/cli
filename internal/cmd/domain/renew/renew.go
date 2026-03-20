package renew

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id          string
	skipConfirm bool
}

func NewCmdRenew(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "renew",
		Short: "Renew a registered domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRenew(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Registered domain ID")
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")

	return cmd
}

func runRenew(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()

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
			options[i] = fmt.Sprintf("%s (expires %s, $%.2f/yr)", d.Domain, d.ExpiresAt.Format("2006-01-02"), float64(d.RenewalPrice)/100)
		}
		idx, err := f.Prompter.Select("Select domain to renew", "", options)
		if err != nil {
			return err
		}
		opts.id = domains[idx].ID
	}

	if f.Interactive && !opts.skipConfirm {
		confirm, err := f.Prompter.Confirm("Renew this domain for 1 year?", false)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Renewing domain..."),
	)
	s.Start()
	domain, err := f.ApiClient.RenewDomain(ctx, opts.id)
	s.Stop()
	if err != nil {
		return fmt.Errorf("renew domain failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(domain)
	}

	f.Log.Infof("Domain %s renewed successfully! New expiry: %s", domain.Domain, domain.ExpiresAt.Format("2006-01-02"))
	return nil
}
