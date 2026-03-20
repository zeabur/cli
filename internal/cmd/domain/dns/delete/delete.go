package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmd/domain/dns/dnsutil"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	domainID    string
	domain      string
	recordID    string
	recordType  string
	name        string
	skipConfirm bool
}

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a DNS record",
		Long:  "Delete a DNS record. Identify the record by --type and --name, or by --record-id.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.domain, "domain", "", "Domain name (e.g. example.com)")
	cmd.Flags().StringVar(&opts.domainID, "domain-id", "", "Registered domain ID (advanced)")
	cmd.Flags().StringVar(&opts.recordType, "type", "", "Record type to match (A, AAAA, CNAME, ...)")
	cmd.Flags().StringVar(&opts.name, "name", "", "Record name to match (e.g. @ or subdomain)")
	cmd.Flags().StringVar(&opts.recordID, "record-id", "", "DNS record ID (advanced)")
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")

	return cmd
}

func runDelete(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()

	if err := resolveDomainID(ctx, f, opts); err != nil {
		return err
	}

	if err := resolveRecordID(ctx, f, opts); err != nil {
		return err
	}

	if f.Interactive && !opts.skipConfirm {
		confirm, err := f.Prompter.Confirm("Delete this DNS record?", false)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	err := f.ApiClient.DeleteDNSRecord(ctx, opts.domainID, opts.recordID)
	if err != nil {
		return fmt.Errorf("delete DNS record failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(map[string]string{"status": "deleted"})
	}

	f.Log.Infof("DNS record deleted")
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

func resolveRecordID(ctx context.Context, f *cmdutil.Factory, opts *Options) error {
	if opts.recordID != "" {
		return nil
	}

	records, err := f.ApiClient.ListDNSRecords(ctx, opts.domainID)
	if err != nil {
		return fmt.Errorf("list DNS records failed: %w", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("no DNS records found")
	}

	if opts.recordType != "" && opts.name != "" {
		matched, err := dnsutil.FindRecord(records, opts.recordType, opts.name)
		if err != nil {
			return err
		}
		if len(matched) == 1 {
			opts.recordID = matched[0].ID
			return nil
		}
		if !f.Interactive {
			return fmt.Errorf("multiple records match type=%s name=%s, use --record-id to specify", opts.recordType, opts.name)
		}
		options := make([]string, len(matched))
		for i, r := range matched {
			options[i] = fmt.Sprintf("%s %s → %s", r.Type, r.Name, r.Content)
		}
		idx, err := f.Prompter.Select("Multiple records found, select one", "", options)
		if err != nil {
			return err
		}
		opts.recordID = matched[idx].ID
		return nil
	}

	if !f.Interactive {
		return fmt.Errorf("--type and --name are required to identify the record (or use --record-id)")
	}

	options := make([]string, len(records))
	for i, r := range records {
		options[i] = fmt.Sprintf("%s %s → %s", r.Type, r.Name, r.Content)
	}
	idx, err := f.Prompter.Select("Select record to delete", "", options)
	if err != nil {
		return err
	}
	opts.recordID = records[idx].ID
	return nil
}
