package dnsutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
)

func ResolveDomainID(ctx context.Context, client api.Client, domain string) (string, error) {
	domains, err := client.ListRegisteredDomains(ctx)
	if err != nil {
		return "", fmt.Errorf("list registered domains failed: %w", err)
	}
	for _, d := range domains {
		if strings.EqualFold(d.Domain, domain) {
			return d.ID, nil
		}
	}
	return "", fmt.Errorf("domain %q not found in your registered domains", domain)
}

func FindRecord(records model.DNSRecords, recordType, name string) ([]model.DNSRecord, error) {
	var matched []model.DNSRecord
	for _, r := range records {
		typeMatch := recordType == "" || strings.EqualFold(r.Type, recordType)
		nameMatch := name == "" || strings.EqualFold(r.Name, name)
		if typeMatch && nameMatch {
			matched = append(matched, r)
		}
	}
	if len(matched) == 0 {
		return nil, fmt.Errorf("no DNS record found matching type=%q name=%q", recordType, name)
	}
	return matched, nil
}
