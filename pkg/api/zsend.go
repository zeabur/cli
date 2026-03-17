package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) GetZSendOnboardingStatus(ctx context.Context) (*model.ZSendOnboardingStatus, error) {
	var query struct {
		GetZSendOnboardingStatus model.ZSendOnboardingStatus `graphql:"getZSendOnboardingStatus"`
	}

	err := c.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return &query.GetZSendOnboardingStatus, nil
}

func (c *client) GetZSendUserStatus(ctx context.Context) (*model.ZSendUserStatus, error) {
	var query struct {
		GetZSendUserStatus model.ZSendUserStatus `graphql:"getZSendUserStatus"`
	}

	err := c.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return &query.GetZSendUserStatus, nil
}

func (c *client) OnboardZSend(ctx context.Context) (*model.ZSendOnboardingStatus, error) {
	var mutation struct {
		OnboardZSend model.ZSendOnboardingStatus `graphql:"onboardZSend"`
	}

	err := c.Mutate(ctx, &mutation, nil)
	if err != nil {
		return nil, err
	}

	return &mutation.OnboardZSend, nil
}

// ListZSendDomains lists ZSend domains with optional pagination.
func (c *client) ListZSendDomains(ctx context.Context, page, pageSize *int) (*model.ListZSendDomainsReply, error) {
	var query struct {
		ListZSendDomains model.ListZSendDomainsReply `graphql:"listZSendDomains(page: $page, pageSize: $pageSize)"`
	}

	err := c.Query(ctx, &query, V{
		"page":     page,
		"pageSize": pageSize,
	})
	if err != nil {
		return nil, err
	}

	return &query.ListZSendDomains, nil
}

// GetZSendDomain returns a ZSend domain by ID.
func (c *client) GetZSendDomain(ctx context.Context, id string) (*model.ZSendDomain, error) {
	var query struct {
		GetZSendDomain model.ZSendDomain `graphql:"getZSendDomain(id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	return &query.GetZSendDomain, nil
}

// CreateZSendDomain creates a new ZSend domain.
func (c *client) CreateZSendDomain(ctx context.Context, domain, region string) (*model.ZSendDomain, error) {
	var mutation struct {
		CreateZSendDomain model.ZSendDomain `graphql:"createZSendDomain(domain: $domain, region: $region)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"domain": domain,
		"region": region,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.CreateZSendDomain, nil
}

// VerifyZSendDomain verifies a ZSend domain.
func (c *client) VerifyZSendDomain(ctx context.Context, id string) (*model.ZSendDomain, error) {
	var mutation struct {
		VerifyZSendDomain model.ZSendDomain `graphql:"verifyZSendDomain(id: $id)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.VerifyZSendDomain, nil
}

// DeleteZSendDomain deletes a ZSend domain.
func (c *client) DeleteZSendDomain(ctx context.Context, id string) error {
	var mutation struct {
		DeleteZSendDomain bool `graphql:"deleteZSendDomain(id: $id)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id": id,
	})
	if err != nil {
		return err
	}

	return nil
}

// ListZSendAPIKeys lists ZSend API keys with optional pagination.
func (c *client) ListZSendAPIKeys(ctx context.Context, page, pageSize *int) (*model.ListZSendAPIKeysReply, error) {
	var query struct {
		ListZSendAPIKeys model.ListZSendAPIKeysReply `graphql:"listZSendAPIKeys(page: $page, pageSize: $pageSize)"`
	}

	err := c.Query(ctx, &query, V{
		"page":     page,
		"pageSize": pageSize,
	})
	if err != nil {
		return nil, err
	}

	return &query.ListZSendAPIKeys, nil
}

// GetZSendAPIKey returns a ZSend API key by ID.
func (c *client) GetZSendAPIKey(ctx context.Context, id string) (*model.ZSendAPIKey, error) {
	var query struct {
		GetZSendAPIKey model.ZSendAPIKey `graphql:"getZSendAPIKey(id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	return &query.GetZSendAPIKey, nil
}

// CreateZSendAPIKey creates a new ZSend API key.
func (c *client) CreateZSendAPIKey(ctx context.Context, input model.CreateZSendAPIKeyInput) (*model.CreateZSendAPIKeyReply, error) {
	var mutation struct {
		CreateZSendAPIKey model.CreateZSendAPIKeyReply `graphql:"createZSendAPIKey(input: $input)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"input": input,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.CreateZSendAPIKey, nil
}

// DeleteZSendAPIKey deletes a ZSend API key.
func (c *client) DeleteZSendAPIKey(ctx context.Context, id string) error {
	var mutation struct {
		DeleteZSendAPIKey bool `graphql:"deleteZSendAPIKey(id: $id)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id": id,
	})
	if err != nil {
		return err
	}

	return nil
}

// ListZSendWebhooks lists ZSend webhooks with optional pagination.
func (c *client) ListZSendWebhooks(ctx context.Context, page, pageSize *int) (*model.ListZSendWebhooksReply, error) {
	var query struct {
		ListZSendWebhooks model.ListZSendWebhooksReply `graphql:"listZSendWebhooks(page: $page, pageSize: $pageSize)"`
	}

	err := c.Query(ctx, &query, V{
		"page":     page,
		"pageSize": pageSize,
	})
	if err != nil {
		return nil, err
	}

	return &query.ListZSendWebhooks, nil
}

// GetZSendWebhook returns a ZSend webhook by ID.
func (c *client) GetZSendWebhook(ctx context.Context, id string) (*model.ZSendWebhook, error) {
	var query struct {
		GetZSendWebhook model.ZSendWebhook `graphql:"getZSendWebhook(id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	return &query.GetZSendWebhook, nil
}

// CreateZSendWebhook creates a new ZSend webhook.
func (c *client) CreateZSendWebhook(ctx context.Context, input model.CreateZSendWebhookInput) (*model.CreateZSendWebhookReply, error) {
	var mutation struct {
		CreateZSendWebhook model.CreateZSendWebhookReply `graphql:"createZSendWebhook(input: $input)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"input": input,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.CreateZSendWebhook, nil
}

// DeleteZSendWebhook deletes a ZSend webhook.
func (c *client) DeleteZSendWebhook(ctx context.Context, id string) error {
	var mutation struct {
		DeleteZSendWebhook bool `graphql:"deleteZSendWebhook(id: $id)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id": id,
	})
	if err != nil {
		return err
	}

	return nil
}

// VerifyZSendWebhook verifies a ZSend webhook.
func (c *client) VerifyZSendWebhook(ctx context.Context, id string) (*model.VerifyZSendWebhookReply, error) {
	var mutation struct {
		VerifyZSendWebhook model.VerifyZSendWebhookReply `graphql:"verifyZSendWebhook(id: $id)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.VerifyZSendWebhook, nil
}

// ListZSendEmails lists ZSend emails with optional pagination and filters.
func (c *client) ListZSendEmails(ctx context.Context, page, pageSize *int, status, jobType, jobID *string) (*model.ListZSendEmailsReply, error) {
	var query struct {
		ListZSendEmails model.ListZSendEmailsReply `graphql:"listZSendEmails(page: $page, pageSize: $pageSize, status: $status, jobType: $jobType, jobId: $jobId)"`
	}

	err := c.Query(ctx, &query, V{
		"page":     page,
		"pageSize": pageSize,
		"status":   status,
		"jobType":  jobType,
		"jobId":    jobID,
	})
	if err != nil {
		return nil, err
	}

	return &query.ListZSendEmails, nil
}

// GetZSendEmail returns a ZSend email by ID.
func (c *client) GetZSendEmail(ctx context.Context, id string) (*model.ZSendEmail, error) {
	var query struct {
		GetZSendEmail model.ZSendEmail `graphql:"getZSendEmail(id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	return &query.GetZSendEmail, nil
}
