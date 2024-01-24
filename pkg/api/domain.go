package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) ListDomains(ctx context.Context, environmentID string) ([]*model.Domain, error) {
	var query struct {
		Service struct {
			Domains []*model.Domain `graphql:"domains(environmentID: $environmentID))"`
		}
	}

	err := c.Query(ctx, &query, V{
		"environmentID": ObjectID(environmentID),
	})

	if err != nil {
		return nil, err
	}

	return query.Service.Domains, nil
}

func (c *client) RemoveDomain(ctx context.Context, domain string) (bool, error) {
	var mutation struct {
		RemoveDomain bool `graphql:"removeDomain(domain: $domain)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"domain": domain,
	})

	if err != nil {
		return false, err
	}

	return mutation.RemoveDomain, nil
}

func (c *client) CheckDomainAvailable(ctx context.Context, domain string, isGenerated bool) (bool, string, error) {
	var mutation struct {
		CheckDomainAvailable struct {
			IsAvailable bool
			Reason      string
		} `graphql:"checkDomainAvailable(domain: $domain, isGenerated: $isGenerated)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"domain":      domain,
		"isGenerated": isGenerated,
	})

	if err != nil {
		return false, "", err
	}

	return mutation.CheckDomainAvailable.IsAvailable, mutation.CheckDomainAvailable.Reason, nil
}
