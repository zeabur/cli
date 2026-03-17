package model

import (
	"fmt"
	"time"
)

// AIHubTenant represents an AI Hub tenant with balance, keys, and auto-recharge settings.
type AIHubTenant struct {
	Balance              int         `graphql:"balance"`
	Keys                 []AIHubKey  `graphql:"keys"`
	ProviderCustomerID   string      `graphql:"providerCustomerID"`
	Provider             string      `graphql:"provider"`
	AutoRechargeThreshold int        `graphql:"autoRechargeThreshold"`
	AutoRechargeAmount   int         `graphql:"autoRechargeAmount"`
}

// AIHubKey represents an API key in AI Hub.
type AIHubKey struct {
	KeyID string  `graphql:"keyID"`
	Alias string  `graphql:"alias"`
	Cost  float64 `graphql:"cost"`
}

// AIHubKeys is a list of AIHubKey for table display.
type AIHubKeys []AIHubKey

func (k AIHubKeys) Header() []string {
	return []string{"Key ID", "Alias", "Cost"}
}

func (k AIHubKeys) Rows() [][]string {
	rows := make([][]string, 0, len(k))
	for _, key := range k {
		rows = append(rows, []string{
			key.KeyID,
			key.Alias,
			fmt.Sprintf("$%.4f", key.Cost),
		})
	}
	return rows
}

var _ Tabler = (AIHubKeys)(nil)

// AddAIHubBalanceResult is the result of adding balance to AI Hub.
type AddAIHubBalanceResult struct {
	NewBalance int `graphql:"newBalance"`
}

// CreateAIHubKeyResult is the result of creating a new AI Hub key.
type CreateAIHubKeyResult struct {
	Key    AIHubKey `graphql:"key"`
	APIKey string   `graphql:"apiKey"`
}

// UpdateAIHubAutoRechargeSettingsResult is the result of updating auto-recharge settings.
type UpdateAIHubAutoRechargeSettingsResult struct {
	AutoRechargeThreshold int `graphql:"autoRechargeThreshold"`
	AutoRechargeAmount    int `graphql:"autoRechargeAmount"`
}

// AIHubSpendLog represents a single spend log entry.
type AIHubSpendLog struct {
	Timestamp        time.Time `graphql:"timestamp"`
	Cost             float64   `graphql:"cost"`
	TotalTokens      int       `graphql:"totalTokens"`
	PromptTokens     int       `graphql:"promptTokens"`
	CompletionTokens int       `graphql:"completionTokens"`
	Model            string    `graphql:"model"`
	KeyAlias         string    `graphql:"keyAlias"`
}

// AIHubSpendLogs is a list of AIHubSpendLog for table display.
type AIHubSpendLogs []AIHubSpendLog

func (l AIHubSpendLogs) Header() []string {
	return []string{"Timestamp", "Model", "Key Alias", "Cost", "Total Tokens", "Prompt Tokens", "Completion Tokens"}
}

func (l AIHubSpendLogs) Rows() [][]string {
	rows := make([][]string, 0, len(l))
	for _, log := range l {
		rows = append(rows, []string{
			log.Timestamp.Format(time.RFC3339),
			log.Model,
			log.KeyAlias,
			fmt.Sprintf("$%.6f", log.Cost),
			fmt.Sprintf("%d", log.TotalTokens),
			fmt.Sprintf("%d", log.PromptTokens),
			fmt.Sprintf("%d", log.CompletionTokens),
		})
	}
	return rows
}

var _ Tabler = (AIHubSpendLogs)(nil)

// AIHubSpendLogsPaginated is a paginated list of spend logs.
type AIHubSpendLogsPaginated struct {
	Data       []AIHubSpendLog `graphql:"data"`
	Total      int             `graphql:"total"`
	Page       int             `graphql:"page"`
	PageSize   int             `graphql:"pageSize"`
	TotalPages int             `graphql:"totalPages"`
}

// AIHubMonthlyUsage represents the monthly usage summary.
type AIHubMonthlyUsage struct {
	TotalSpend float64          `graphql:"totalSpend"`
	DailyUsage []AIHubDailyUsage `graphql:"dailyUsage"`
	ModelsCost []AIHubModelCost  `graphql:"modelsCost"`
}

// AIHubDailyUsage represents a single day's usage.
type AIHubDailyUsage struct {
	Date   string           `graphql:"date"`
	Spend  float64          `graphql:"spend"`
	Models []AIHubModelCost `graphql:"models"`
}

// AIHubModelCost represents the cost for a specific model.
type AIHubModelCost struct {
	Model string  `graphql:"model"`
	Cost  float64 `graphql:"cost"`
}

// AIHubModelCosts is a list of AIHubModelCost for table display.
type AIHubModelCosts []AIHubModelCost

func (m AIHubModelCosts) Header() []string {
	return []string{"Model", "Cost"}
}

func (m AIHubModelCosts) Rows() [][]string {
	rows := make([][]string, 0, len(m))
	for _, mc := range m {
		rows = append(rows, []string{
			mc.Model,
			fmt.Sprintf("$%.6f", mc.Cost),
		})
	}
	return rows
}

var _ Tabler = (AIHubModelCosts)(nil)
