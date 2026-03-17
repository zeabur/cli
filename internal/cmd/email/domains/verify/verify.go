package verify

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string
}

func NewCmdVerify(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify an email domain's DNS records",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVerify(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Domain ID")

	return cmd
}

func runVerify(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runVerifyInteractive(f, opts)
	}
	return runVerifyNonInteractive(f, opts)
}

func runVerifyInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.id == "" {
		id, err := f.Prompter.Input("Domain ID: ", "")
		if err != nil {
			return err
		}
		opts.id = id
	}

	if err := paramCheck(opts); err != nil {
		return err
	}
	return verifyDomain(f, opts)
}

func runVerifyNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return verifyDomain(f, opts)
}

func verifyDomain(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Verifying domain..."),
	)
	s.Start()
	domain, err := f.ApiClient.VerifyZSendDomain(context.Background(), opts.id)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(domain)
	}

	f.Log.Infof("Domain %q verification triggered (Status: %s)", domain.Value, domain.Status)

	return nil
}

func paramCheck(opts Options) error {
	if opts.id == "" {
		return fmt.Errorf("domain ID is required")
	}
	return nil
}
