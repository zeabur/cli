package create

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
	recordType string
	name       string
	content    string
	ttl        int
	priority   int
	proxied    bool
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a DNS record",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.domain, "domain", "", "Domain name (e.g. example.com)")
	cmd.Flags().StringVar(&opts.domainID, "domain-id", "", "Registered domain ID (advanced)")
	cmd.Flags().StringVar(&opts.recordType, "type", "", "Record type (A, AAAA, CNAME, MX, TXT, SRV, CAA, NS)")
	cmd.Flags().StringVar(&opts.name, "name", "", "Record name (e.g. @ or subdomain)")
	cmd.Flags().StringVar(&opts.content, "content", "", "Record content")
	cmd.Flags().IntVar(&opts.ttl, "ttl", 1, "TTL (1 = auto)")
	cmd.Flags().IntVar(&opts.priority, "priority", 0, "Priority (for MX, SRV)")
	cmd.Flags().BoolVar(&opts.proxied, "proxied", false, "Proxied through Cloudflare")

	return cmd
}

func runCreate(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()

	if err := resolveDomainID(ctx, f, opts); err != nil {
		return err
	}

	if f.Interactive {
		if opts.recordType == "" {
			types := []string{"A", "AAAA", "CNAME", "MX", "TXT", "SRV", "CAA", "NS"}
			idx, err := f.Prompter.Select("Record type", "", types)
			if err != nil {
				return err
			}
			opts.recordType = types[idx]
		}
		if opts.name == "" {
			name, err := f.Prompter.Input("Record name (e.g. @ or subdomain): ", "")
			if err != nil {
				return err
			}
			opts.name = name
		}
		if opts.content == "" {
			content, err := f.Prompter.Input("Record content: ", "")
			if err != nil {
				return err
			}
			opts.content = content
		}
	}

	if opts.recordType == "" {
		return fmt.Errorf("--type is required")
	}
	if opts.name == "" {
		return fmt.Errorf("--name is required")
	}
	if opts.content == "" {
		return fmt.Errorf("--content is required")
	}

	input := model.CreateDNSRecordInput{
		Type:    model.RegisteredDomainDNSRecordType(opts.recordType),
		Name:    opts.name,
		Content: opts.content,
	}
	if opts.ttl != 0 {
		input.TTL = &opts.ttl
	}
	if opts.priority != 0 {
		input.Priority = &opts.priority
	}
	if opts.proxied {
		input.Proxied = &opts.proxied
	}

	record, err := f.ApiClient.CreateDNSRecord(ctx, opts.domainID, input)
	if err != nil {
		return fmt.Errorf("create DNS record failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(record)
	}

	f.Log.Infof("DNS record created: %s %s %s", record.Type, record.Name, record.Content)
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
