package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) AddDomain(ctx context.Context, serviceID string, environmentID string, isGenerated bool, domain string, options ...string) (*string, error) {
	var err error

	if len(options) > 0 {
		var mutationOptional struct {
			AddDomain struct {
				Domain string `json:"domain" graphql:"domain"`
			} `graphql:"addDomain(serviceID: $serviceID, environmentID: $environmentID, isGenerated: $isGenerated, domain: $domain, redirectTo: $redirectTo)"`
		}

		err = c.Mutate(ctx, &mutationOptional, V{
			"serviceID":     ObjectID(serviceID),
			"environmentID": ObjectID(environmentID),
			"isGenerated":   isGenerated,
			"domain":        domain,
			"redirectTo":    options[0],
		})
		if err != nil {
			return nil, err
		}

		return &mutationOptional.AddDomain.Domain, nil
	}

	var mutation struct {
		AddDomain struct {
			Domain string `json:"domain" graphql:"domain"`
		} `graphql:"addDomain(serviceID: $serviceID, environmentID: $environmentID, isGenerated: $isGenerated, domain: $domain)"`
	}

	err = c.Mutate(ctx, &mutation, V{
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
		"isGenerated":   isGenerated,
		"domain":        domain,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.AddDomain.Domain, nil
}

func (c *client) ListDomains(ctx context.Context, serviceID string, environmentID string) (model.Domains, error) {
	var query struct {
		Service struct {
			Domains model.Domains `graphql:"domains(environmentID: $environmentID)"`
		} `graphql:"service(_id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"environmentID": ObjectID(environmentID),
		"id":            ObjectID(serviceID),
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

func (c *client) CheckDomainAvailable(ctx context.Context, domain string, isGenerated bool, region string) (bool, string, error) {
	var mutation struct {
		CheckDomainAvailable struct {
			IsAvailable bool
			Reason      string
		} `graphql:"checkDomainAvailable(domain: $domain, isGenerated: $isGenerated, region: $region)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"domain":      domain,
		"isGenerated": isGenerated,
		"region":      region,
	})
	if err != nil {
		return false, "", err
	}

	return mutation.CheckDomainAvailable.IsAvailable, mutation.CheckDomainAvailable.Reason, nil
}
