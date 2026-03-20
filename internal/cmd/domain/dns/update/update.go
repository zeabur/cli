package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmd/domain/dns/dnsutil"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	domainID   string
	domain     string
	recordID   string
	recordType string
	name       string
	content    string
	ttl        int
	priority   int
	proxied    bool
}

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	var ttlSet, prioritySet, proxiedSet bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a DNS record",
		Long:  "Update a DNS record. Identify the record by --type and --name, or by --record-id.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ttlSet = cmd.Flags().Changed("ttl")
			prioritySet = cmd.Flags().Changed("priority")
			proxiedSet = cmd.Flags().Changed("proxied")
			return runUpdate(f, opts, ttlSet, prioritySet, proxiedSet)
		},
	}

	cmd.Flags().StringVar(&opts.domain, "domain", "", "Domain name (e.g. example.com)")
	cmd.Flags().StringVar(&opts.domainID, "domain-id", "", "Registered domain ID (advanced)")
	cmd.Flags().StringVar(&opts.recordType, "type", "", "Record type to match (A, AAAA, CNAME, ...)")
	cmd.Flags().StringVar(&opts.name, "name", "", "Record name to match (e.g. @ or subdomain)")
	cmd.Flags().StringVar(&opts.recordID, "record-id", "", "DNS record ID (advanced)")
	cmd.Flags().StringVar(&opts.content, "content", "", "New record content")
	cmd.Flags().IntVar(&opts.ttl, "ttl", 0, "New TTL")
	cmd.Flags().IntVar(&opts.priority, "priority", 0, "New priority")
	cmd.Flags().BoolVar(&opts.proxied, "proxied", false, "Proxied through Cloudflare")

	return cmd
}

func runUpdate(f *cmdutil.Factory, opts *Options, ttlSet, prioritySet, proxiedSet bool) error {
	ctx := context.Background()

	if err := resolveDomainID(ctx, f, opts); err != nil {
		return err
	}

	if err := resolveRecordID(ctx, f, opts); err != nil {
		return err
	}

	if f.Interactive && opts.content == "" {
		content, err := f.Prompter.Input("New content (leave empty to skip): ", "")
		if err != nil {
			return err
		}
		opts.content = content
	}

	input := model.UpdateDNSRecordInput{}
	if opts.content != "" {
		input.Content = &opts.content
	}
	if ttlSet {
		input.TTL = &opts.ttl
	}
	if prioritySet {
		input.Priority = &opts.priority
	}
	if proxiedSet {
		input.Proxied = &opts.proxied
	}

	record, err := f.ApiClient.UpdateDNSRecord(ctx, opts.domainID, opts.recordID, input)
	if err != nil {
		return fmt.Errorf("update DNS record failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(record)
	}

	f.Log.Infof("DNS record updated: %s %s %s", record.Type, record.Name, record.Content)
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

	// If type+name provided, find matching record
	if opts.recordType != "" && opts.name != "" {
		matched, err := dnsutil.FindRecord(records, opts.recordType, opts.name)
		if err != nil {
			return err
		}
		if len(matched) == 1 {
			opts.recordID = matched[0].ID
			return nil
		}
		// Multiple matches — let user pick in interactive mode
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

	// Interactive: let user pick
	if !f.Interactive {
		return fmt.Errorf("--type and --name are required to identify the record (or use --record-id)")
	}

	options := make([]string, len(records))
	for i, r := range records {
		options[i] = fmt.Sprintf("%s %s → %s", r.Type, r.Name, r.Content)
	}
	idx, err := f.Prompter.Select("Select record to update", "", options)
	if err != nil {
		return err
	}
	opts.recordID = records[idx].ID
	return nil
}
