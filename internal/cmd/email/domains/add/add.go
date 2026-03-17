package add

import (
	"context"
	"fmt"
	"slices"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

var regionChoices = []string{
	"us-east-1",
	"us-west-1",
	"eu-central-1",
	"ap-northeast-1",
	"ap-northeast-3",
	"ap-southeast-1",
}

type Options struct {
	domain string
	region string
}

func NewCmdAdd(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an email domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.domain, "domain", "", "Domain name")
	cmd.Flags().StringVar(&opts.region, "region", "", "Region (e.g. us-east-1)")

	return cmd
}

func runAdd(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runAddInteractive(f, opts)
	}
	return runAddNonInteractive(f, opts)
}

func runAddInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.domain == "" {
		domain, err := f.Prompter.Input("Domain: ", "")
		if err != nil {
			return err
		}
		opts.domain = domain
	}

	if opts.region == "" {
		idx, err := f.Prompter.Select("Region: ", "", regionChoices)
		if err != nil {
			return err
		}
		opts.region = regionChoices[idx]
	}

	if err := paramCheck(opts); err != nil {
		return err
	}

	return addDomain(f, opts)
}

func runAddNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return addDomain(f, opts)
}

func addDomain(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Adding domain..."),
	)
	s.Start()
	domain, err := f.ApiClient.CreateZSendDomain(context.Background(), opts.domain, opts.region)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(domain)
	}

	f.Log.Infof("Domain %q added successfully (ID: %s, Status: %s)", domain.Value, domain.ID, domain.Status)

	return nil
}

func paramCheck(opts Options) error {
	if opts.domain == "" {
		return fmt.Errorf("domain is required")
	}
	if opts.region == "" {
		return fmt.Errorf("region is required")
	}
	if !slices.Contains(regionChoices, opts.region) {
		return fmt.Errorf("unsupported region %q", opts.region)
	}
	return nil
}
