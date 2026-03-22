package verification

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type statusOptions struct {
	id string
}

func newCmdStatus(f *cmdutil.Factory) *cobra.Command {
	opts := &statusOptions{}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check the ICANN registrant verification status of a domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Registered domain ID")

	return cmd
}

func runStatus(f *cmdutil.Factory, opts *statusOptions) error {
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
			options[i] = d.Domain
		}
		idx, err := f.Prompter.Select("Select domain to check verification status", "", options)
		if err != nil {
			return err
		}
		opts.id = domains[idx].ID
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Checking verification status..."),
	)
	s.Start()
	domain, err := f.ApiClient.GetRegisteredDomain(ctx, opts.id)
	s.Stop()
	if err != nil {
		return fmt.Errorf("get registered domain failed: %w", err)
	}

	status := "N/A"
	if domain.RegistrantVerificationStatus != nil {
		status = *domain.RegistrantVerificationStatus
	}

	if f.JSON {
		return f.Printer.JSON(map[string]string{
			"domain": domain.Domain,
			"status": status,
		})
	}

	f.Printer.Table(
		[]string{"Domain", "Verification Status"},
		[][]string{{domain.Domain, status}},
	)

	return nil
}
