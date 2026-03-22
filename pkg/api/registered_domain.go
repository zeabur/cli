package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) CheckDomainRegistrationAvailability(ctx context.Context, domain string) (*model.DomainSearchResult, error) {
	var q struct {
		Result model.DomainSearchResult `graphql:"checkDomainRegistrationAvailability(domain: $domain)"`
	}

	err := c.Query(ctx, &q, V{
		"domain": domain,
	})
	if err != nil {
		return nil, err
	}

	return &q.Result, nil
}

func (c *client) PurchaseDomain(ctx context.Context, domain, registrantProfileID string) (*model.PurchaseDomainResult, error) {
	var mutation struct {
		PurchaseDomain model.PurchaseDomainResult `graphql:"purchaseDomain(domain: $domain, registrantProfileID: $registrantProfileID)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"domain":              domain,
		"registrantProfileID": ObjectID(registrantProfileID),
	})
	if err != nil {
		return nil, err
	}

	return &mutation.PurchaseDomain, nil
}

func (c *client) ListRegisteredDomains(ctx context.Context) (model.RegisteredDomains, error) {
	var query struct {
		RegisteredDomains model.RegisteredDomains `graphql:"registeredDomains"`
	}

	err := c.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return query.RegisteredDomains, nil
}

func (c *client) GetRegisteredDomain(ctx context.Context, id string) (*model.RegisteredDomain, error) {
	var query struct {
		RegisteredDomain model.RegisteredDomain `graphql:"registeredDomain(id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id": ObjectID(id),
	})
	if err != nil {
		return nil, err
	}

	return &query.RegisteredDomain, nil
}

func (c *client) RenewDomain(ctx context.Context, id string) (*model.RegisteredDomain, error) {
	var mutation struct {
		RenewDomain model.RegisteredDomain `graphql:"renewDomain(id: $id)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id": ObjectID(id),
	})
	if err != nil {
		return nil, err
	}

	return &mutation.RenewDomain, nil
}

func (c *client) SetDomainAutoRenew(ctx context.Context, id string, autoRenew bool) (*model.RegisteredDomain, error) {
	var mutation struct {
		SetDomainAutoRenew model.RegisteredDomain `graphql:"setDomainAutoRenew(id: $id, autoRenew: $autoRenew)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id":        ObjectID(id),
		"autoRenew": autoRenew,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.SetDomainAutoRenew, nil
}

func (c *client) ListDNSRecords(ctx context.Context, registeredDomainID string) (model.DNSRecords, error) {
	var query struct {
		DNSRecords model.DNSRecords `graphql:"registeredDomainDNSRecords(registeredDomainID: $registeredDomainID)"`
	}

	err := c.Query(ctx, &query, V{
		"registeredDomainID": ObjectID(registeredDomainID),
	})
	if err != nil {
		return nil, err
	}

	return query.DNSRecords, nil
}

func (c *client) CreateDNSRecord(ctx context.Context, registeredDomainID string, input model.CreateDNSRecordInput) (*model.DNSRecord, error) {
	var mutation struct {
		CreateDNSRecord model.DNSRecord `graphql:"createRegisteredDomainDNSRecord(registeredDomainID: $registeredDomainID, input: $input)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"registeredDomainID": ObjectID(registeredDomainID),
		"input":              input,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.CreateDNSRecord, nil
}

func (c *client) UpdateDNSRecord(ctx context.Context, registeredDomainID, recordID string, input model.UpdateDNSRecordInput) (*model.DNSRecord, error) {
	var mutation struct {
		UpdateDNSRecord model.DNSRecord `graphql:"updateRegisteredDomainDNSRecord(registeredDomainID: $registeredDomainID, recordID: $recordID, input: $input)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"registeredDomainID": ObjectID(registeredDomainID),
		"recordID":           recordID,
		"input":              input,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.UpdateDNSRecord, nil
}

func (c *client) DeleteDNSRecord(ctx context.Context, registeredDomainID, recordID string) error {
	var mutation struct {
		DeleteDNSRecord bool `graphql:"deleteRegisteredDomainDNSRecord(registeredDomainID: $registeredDomainID, recordID: $recordID)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"registeredDomainID": ObjectID(registeredDomainID),
		"recordID":           recordID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) ListRegistrantProfiles(ctx context.Context) (model.RegistrantProfiles, error) {
	var query struct {
		RegistrantProfiles model.RegistrantProfiles `graphql:"registrantProfiles"`
	}

	err := c.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return query.RegistrantProfiles, nil
}

func (c *client) CreateRegistrantProfile(ctx context.Context, input model.CreateRegistrantProfileInput) (*model.RegistrantProfile, error) {
	var mutation struct {
		CreateRegistrantProfile model.RegistrantProfile `graphql:"createRegistrantProfile(input: $input)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"input": input,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.CreateRegistrantProfile, nil
}

func (c *client) UpdateRegistrantProfile(ctx context.Context, id string, input model.UpdateRegistrantProfileInput) (*model.RegistrantProfile, error) {
	var mutation struct {
		UpdateRegistrantProfile model.RegistrantProfile `graphql:"updateRegistrantProfile(id: $id, input: $input)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id":    ObjectID(id),
		"input": input,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.UpdateRegistrantProfile, nil
}

func (c *client) ResendRegistrantVerificationEmail(ctx context.Context, registeredDomainID string) error {
	var mutation struct {
		ResendRegistrantVerificationEmail bool `graphql:"resendRegistrantVerificationEmail(registeredDomainID: $registeredDomainID)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"registeredDomainID": ObjectID(registeredDomainID),
	})

	return err
}

func (c *client) UpdateRegistrantContact(ctx context.Context, registeredDomainID string, input model.UpdateRegistrantContactInput) error {
	var mutation struct {
		UpdateRegistrantContact bool `graphql:"updateRegistrantContact(registeredDomainID: $registeredDomainID, input: $input)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"registeredDomainID": ObjectID(registeredDomainID),
		"input":              input,
	})

	return err
}

func (c *client) DeleteRegistrantProfile(ctx context.Context, id string) error {
	var mutation struct {
		DeleteRegistrantProfile bool `graphql:"deleteRegistrantProfile(id: $id)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id": ObjectID(id),
	})
	if err != nil {
		return err
	}

	return nil
}
