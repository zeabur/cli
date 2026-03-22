package verification

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdVerification(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verification",
		Short: "Manage registrant verification for registered domains",
	}

	cmd.AddCommand(newCmdStatus(f))
	cmd.AddCommand(newCmdResend(f))
	cmd.AddCommand(newCmdUpdateContact(f))

	return cmd
}

type resendOptions struct {
	id string
}

func newCmdResend(f *cmdutil.Factory) *cobra.Command {
	opts := &resendOptions{}

	cmd := &cobra.Command{
		Use:   "resend",
		Short: "Resend the ICANN registrant verification email",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runResend(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Registered domain ID")

	return cmd
}

func runResend(f *cmdutil.Factory, opts *resendOptions) error {
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
		idx, err := f.Prompter.Select("Select domain to resend verification email", "", options)
		if err != nil {
			return err
		}
		opts.id = domains[idx].ID
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Resending verification email..."),
	)
	s.Start()
	err := f.ApiClient.ResendRegistrantVerificationEmail(ctx, opts.id)
	s.Stop()
	if err != nil {
		return fmt.Errorf("resend verification email failed: %w", err)
	}

	f.Log.Infof("Verification email resent successfully")
	return nil
}
