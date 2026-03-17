package model

import (
	"fmt"
	"time"

	"github.com/zeabur/cli/pkg/util"
)

// ZSendOnboardingStatus represents the onboarding status of ZSend.
type ZSendOnboardingStatus struct {
	IsNew     bool `graphql:"isNew"`
	Submitted bool `graphql:"submitted"`
}

// ZSendUserStatus represents the user status of ZSend.
type ZSendUserStatus struct {
	ID                string     `graphql:"id"`
	UserID            string     `graphql:"userID"`
	Status            string     `graphql:"status"`
	StatusMsg         *string    `graphql:"statusMsg"`
	DailyQuota        int        `graphql:"dailyQuota"`
	DailySent         int        `graphql:"dailySent"`
	QuotaResetAt      time.Time  `graphql:"quotaResetAt"`
	MonthlyQuota      int        `graphql:"monthlyQuota"`
	MonthlySent       int        `graphql:"monthlySent"`
	MonthlyResetAt    *time.Time `graphql:"monthlyResetAt"`
	QuotaType         string     `graphql:"quotaType"`
	QuotaMode         string     `graphql:"quotaMode"`
	ResourceLimitMode string     `graphql:"resourceLimitMode"`
	MaxDomains        int        `graphql:"maxDomains"`
	MaxAPIKeys        int        `graphql:"maxAPIKeys"`
	MaxWebhooks       int        `graphql:"maxWebhooks"`
	SentCount24h      int        `graphql:"sentCount24h"`
	BounceCount24h    int        `graphql:"bounceCount24h"`
	ComplaintCount24h int        `graphql:"complaintCount24h"`
	LastRiskCheckAt   time.Time  `graphql:"lastRiskCheckAt"`
	CreatedAt         time.Time  `graphql:"createdAt"`
	UpdatedAt         time.Time  `graphql:"updatedAt"`
}

// ZSendDNSRecord represents a DNS record for a ZSend domain.
type ZSendDNSRecord struct {
	Category      string     `graphql:"category"`
	Type          string     `graphql:"type"`
	Name          string     `graphql:"name"`
	Content       string     `graphql:"content"`
	TTL           string     `graphql:"ttl"`
	Priority      string     `graphql:"priority"`
	Status        string     `graphql:"status"`
	LastCheckedAt *time.Time `graphql:"lastCheckedAt"`
}

// ZSendDomain represents a ZSend email domain.
type ZSendDomain struct {
	ID        string           `graphql:"id"`
	UserID    string           `graphql:"userID"`
	Type      string           `graphql:"type"`
	Value     string           `graphql:"value"`
	Region    string           `graphql:"region"`
	Status    string           `graphql:"status"`
	StatusMsg *string          `graphql:"statusMsg"`
	Records   []ZSendDNSRecord `graphql:"records"`
	CreatedAt time.Time        `graphql:"createdAt"`
	UpdatedAt time.Time        `graphql:"updatedAt"`
}

// ListZSendDomainsReply is the response for listing ZSend domains.
type ListZSendDomainsReply struct {
	Domains    []ZSendDomain `graphql:"domains"`
	TotalCount int           `graphql:"totalCount"`
}

// ZSendDomains is a list of ZSendDomain for table display.
type ZSendDomains []*ZSendDomain

func (d ZSendDomains) Header() []string {
	return []string{"ID", "Domain", "Region", "Status", "Created At"}
}

func (d ZSendDomains) Rows() [][]string {
	rows := make([][]string, 0, len(d))
	for _, item := range d {
		rows = append(rows, []string{
			item.ID,
			item.Value,
			item.Region,
			item.Status,
			util.ConvertTimeAgoString(item.CreatedAt),
		})
	}
	return rows
}

var _ Tabler = (ZSendDomains)(nil)

// ZSendAPIKey represents a ZSend API key.
type ZSendAPIKey struct {
	ID         string     `graphql:"id"`
	UserID     string     `graphql:"userID"`
	StatusID   string     `graphql:"statusID"`
	Name       string     `graphql:"name"`
	Permission string     `graphql:"permission"`
	Domains    []string   `graphql:"domains"`
	Token      *string    `graphql:"token"`
	CreatedAt  time.Time  `graphql:"createdAt"`
	LastUsedAt *time.Time `graphql:"lastUsedAt"`
}

// ListZSendAPIKeysReply is the response for listing ZSend API keys.
type ListZSendAPIKeysReply struct {
	APIKeys    []ZSendAPIKey `graphql:"apiKeys"`
	TotalCount int           `graphql:"totalCount"`
}

// CreateZSendAPIKeyInput is the input for creating a ZSend API key.
type CreateZSendAPIKeyInput struct {
	Name       string   `json:"name" graphql:"name"`
	Permission string   `json:"permission" graphql:"permission"`
	Domains    []string `json:"domains,omitempty" graphql:"domains"`
}

// CreateZSendAPIKeyReply is the response for creating a ZSend API key.
type CreateZSendAPIKeyReply struct {
	APIKey ZSendAPIKey `graphql:"apiKey"`
}

// UpdateZSendAPIKeyInput is the input for updating a ZSend API key.
type UpdateZSendAPIKeyInput struct {
	ID         string   `json:"id" graphql:"id"`
	Name       *string  `json:"name,omitempty" graphql:"name"`
	Permission *string  `json:"permission,omitempty" graphql:"permission"`
	Domains    []string `json:"domains,omitempty" graphql:"domains"`
}

// ZSendAPIKeys is a list of ZSendAPIKey for table display.
type ZSendAPIKeys []*ZSendAPIKey

func (k ZSendAPIKeys) Header() []string {
	return []string{"ID", "Name", "Permission", "Created At"}
}

func (k ZSendAPIKeys) Rows() [][]string {
	rows := make([][]string, 0, len(k))
	for _, item := range k {
		rows = append(rows, []string{
			item.ID,
			item.Name,
			item.Permission,
			util.ConvertTimeAgoString(item.CreatedAt),
		})
	}
	return rows
}

var _ Tabler = (ZSendAPIKeys)(nil)

// ZSendWebhook represents a ZSend webhook.
type ZSendWebhook struct {
	ID           string     `graphql:"id"`
	UserID       string     `graphql:"userID"`
	StatusID     string     `graphql:"statusID"`
	Name         string     `graphql:"name"`
	Endpoint     string     `graphql:"endpoint"`
	Events       []string   `graphql:"events"`
	Status       string     `graphql:"status"`
	StatusMsg    *string    `graphql:"statusMsg"`
	Enabled      bool       `graphql:"enabled"`
	TotalSent    int        `graphql:"totalSent"`
	SuccessCount int        `graphql:"successCount"`
	FailureCount int        `graphql:"failureCount"`
	LastSentAt   *time.Time `graphql:"lastSentAt"`
	LastError    *string    `graphql:"lastError"`
	CreatedAt    time.Time  `graphql:"createdAt"`
	UpdatedAt    time.Time  `graphql:"updatedAt"`
}

// ListZSendWebhooksReply is the response for listing ZSend webhooks.
type ListZSendWebhooksReply struct {
	Webhooks   []ZSendWebhook `graphql:"webhooks"`
	TotalCount int            `graphql:"totalCount"`
}

// CreateZSendWebhookInput is the input for creating a ZSend webhook.
type CreateZSendWebhookInput struct {
	Name     string   `json:"name" graphql:"name"`
	Endpoint string   `json:"endpoint" graphql:"endpoint"`
	Events   []string `json:"events,omitempty" graphql:"events"`
	Enabled  *bool    `json:"enabled,omitempty" graphql:"enabled"`
}

// CreateZSendWebhookReply is the response for creating a ZSend webhook.
type CreateZSendWebhookReply struct {
	Webhook ZSendWebhook `graphql:"webhook"`
	Secret  string       `graphql:"secret"`
}

// UpdateZSendWebhookInput is the input for updating a ZSend webhook.
type UpdateZSendWebhookInput struct {
	ID       string   `json:"id" graphql:"id"`
	Name     *string  `json:"name,omitempty" graphql:"name"`
	Endpoint *string  `json:"endpoint,omitempty" graphql:"endpoint"`
	Events   []string `json:"events,omitempty" graphql:"events"`
	Enabled  *bool    `json:"enabled,omitempty" graphql:"enabled"`
}

// VerifyZSendWebhookReply is the response for verifying a ZSend webhook.
type VerifyZSendWebhookReply struct {
	Success bool   `graphql:"success"`
	Message string `graphql:"message"`
}

// ZSendWebhooks is a list of ZSendWebhook for table display.
type ZSendWebhooks []*ZSendWebhook

func (w ZSendWebhooks) Header() []string {
	return []string{"ID", "Name", "Endpoint", "Status", "Enabled", "Created At"}
}

func (w ZSendWebhooks) Rows() [][]string {
	rows := make([][]string, 0, len(w))
	for _, item := range w {
		enabled := "No"
		if item.Enabled {
			enabled = "Yes"
		}
		rows = append(rows, []string{
			item.ID,
			item.Name,
			item.Endpoint,
			item.Status,
			enabled,
			util.ConvertTimeAgoString(item.CreatedAt),
		})
	}
	return rows
}

var _ Tabler = (ZSendWebhooks)(nil)

// ZSendEmail represents a ZSend email.
type ZSendEmail struct {
	ID          string     `graphql:"id"`
	UserID      string     `graphql:"userID"`
	APIKeyID    string     `graphql:"apiKeyID"`
	JobType     string     `graphql:"jobType"`
	JobID       string     `graphql:"jobID"`
	MessageID   string     `graphql:"messageID"`
	From        string     `graphql:"from"`
	To          []string   `graphql:"to"`
	CC          []string   `graphql:"cc"`
	BCC         []string   `graphql:"bcc"`
	ReplyTo     []string   `graphql:"replyTo"`
	Subject     string     `graphql:"subject"`
	HTML        string     `graphql:"html"`
	Text        string     `graphql:"text"`
	Status      string     `graphql:"status"`
	MalStatus   string     `graphql:"malStatus"`
	CreatedAt   time.Time  `graphql:"createdAt"`
	ScheduledAt *time.Time `graphql:"scheduledAt"`
}

// ListZSendEmailsReply is the response for listing ZSend emails.
type ListZSendEmailsReply struct {
	Emails     []ZSendEmail `graphql:"emails"`
	TotalCount int          `graphql:"totalCount"`
}

// ZSendEmails is a list of ZSendEmail for table display.
type ZSendEmails []*ZSendEmail

func (e ZSendEmails) Header() []string {
	return []string{"ID", "From", "To", "Subject", "Status", "Created At"}
}

func (e ZSendEmails) Rows() [][]string {
	rows := make([][]string, 0, len(e))
	for _, item := range e {
		to := ""
		if len(item.To) > 0 {
			to = item.To[0]
			if len(item.To) > 1 {
				to = fmt.Sprintf("%s (+%d)", to, len(item.To)-1)
			}
		}
		rows = append(rows, []string{
			item.ID,
			item.From,
			to,
			item.Subject,
			item.Status,
			util.ConvertTimeAgoString(item.CreatedAt),
		})
	}
	return rows
}

var _ Tabler = (ZSendEmails)(nil)
