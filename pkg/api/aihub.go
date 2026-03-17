package api

import (
	"context"
	"time"

	"github.com/zeabur/cli/pkg/model"
)

// GetAIHubTenant returns the current user's AI Hub tenant information.
func (c *client) GetAIHubTenant(ctx context.Context) (*model.AIHubTenant, error) {
	var query struct {
		AIHubTenant model.AIHubTenant `graphql:"aihubTenant"`
	}

	err := c.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return &query.AIHubTenant, nil
}

// AddAIHubBalance adds balance to the AI Hub account.
func (c *client) AddAIHubBalance(ctx context.Context, amount int, provider *string) (*model.AddAIHubBalanceResult, error) {
	var mutation struct {
		AddAIHubBalance model.AddAIHubBalanceResult `graphql:"addAIHubBalance(amount: $amount, provider: $provider)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"amount":   amount,
		"provider": provider,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.AddAIHubBalance, nil
}

// CreateAIHubKey creates a new AI Hub API key.
func (c *client) CreateAIHubKey(ctx context.Context, alias *string) (*model.CreateAIHubKeyResult, error) {
	var mutation struct {
		CreateAIHubKey model.CreateAIHubKeyResult `graphql:"createAIHubKey(alias: $alias)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"alias": alias,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.CreateAIHubKey, nil
}

// DeleteAIHubKey deletes an AI Hub API key by its key ID.
func (c *client) DeleteAIHubKey(ctx context.Context, keyID string) error {
	var mutation struct {
		DeleteAIHubKey bool `graphql:"deleteAIHubKey(keyID: $keyID)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"keyID": keyID,
	})
	if err != nil {
		return err
	}

	return nil
}

// UpdateAIHubAutoRechargeSettings updates the auto-recharge settings.
func (c *client) UpdateAIHubAutoRechargeSettings(ctx context.Context, threshold, amount int) (*model.UpdateAIHubAutoRechargeSettingsResult, error) {
	var mutation struct {
		UpdateAIHubAutoRechargeSettings model.UpdateAIHubAutoRechargeSettingsResult `graphql:"updateAIHubAutoRechargeSettings(threshold: $threshold, amount: $amount)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"threshold": threshold,
		"amount":    amount,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.UpdateAIHubAutoRechargeSettings, nil
}

// GetAIHubSpendLogs returns spend logs optionally filtered by date range.
func (c *client) GetAIHubSpendLogs(ctx context.Context, startDate, endDate *time.Time) ([]model.AIHubSpendLog, error) {
	var query struct {
		AIHubSpendLogs []model.AIHubSpendLog `graphql:"aihubSpendLogs(startDate: $startDate, endDate: $endDate)"`
	}

	err := c.Query(ctx, &query, V{
		"startDate": startDate,
		"endDate":   endDate,
	})
	if err != nil {
		return nil, err
	}

	return query.AIHubSpendLogs, nil
}

// GetAIHubMonthlyUsage returns monthly usage summary for the given month.
func (c *client) GetAIHubMonthlyUsage(ctx context.Context, month *string) (*model.AIHubMonthlyUsage, error) {
	var query struct {
		AIHubMonthlyUsage model.AIHubMonthlyUsage `graphql:"aihubMonthlyUsage(month: $month)"`
	}

	err := c.Query(ctx, &query, V{
		"month": month,
	})
	if err != nil {
		return nil, err
	}

	return &query.AIHubMonthlyUsage, nil
}
