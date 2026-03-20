package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmd/domain/dns/dnsutil"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	domainID string
	domain   string
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List DNS records",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.domain, "domain", "", "Domain name (e.g. example.com)")
	cmd.Flags().StringVar(&opts.domainID, "domain-id", "", "Registered domain ID (advanced)")

	return cmd
}

func runList(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()

	if err := resolveDomainID(ctx, f, opts); err != nil {
		return err
	}

	records, err := f.ApiClient.ListDNSRecords(ctx, opts.domainID)
	if err != nil {
		return fmt.Errorf("list DNS records failed: %w", err)
	}

	if len(records) == 0 {
		if f.JSON {
			return f.Printer.JSON([]any{})
		}
		f.Log.Infof("No DNS records found")
		return nil
	}

	if f.JSON {
		return f.Printer.JSON(records)
	}

	f.Printer.Table(records.Header(), records.Rows())
	return nil
}

func resolveDomainID(ctx context.Context, f *cmdutil.Factory, opts *Options) error {
	if opts.domainID != "" {
		return nil
	}

	if opts.domain != "" {
		id, err := dnsutil.ResolveDomainID(ctx, f.ApiClient, opts.domain)
		if err != nil {
			return err
		}
		opts.domainID = id
		return nil
	}

	if !f.Interactive {
		return fmt.Errorf("--domain is required")
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
	idx, err := f.Prompter.Select("Select domain", "", options)
	if err != nil {
		return err
	}
	opts.domainID = domains[idx].ID
	return nil
}
