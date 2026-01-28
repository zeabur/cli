package model

import (
	"time"

	"github.com/zeabur/cli/pkg/util"
)

// BannedReasonType is the type of reason a user is banned
type BannedReasonType string

// valid banned reasons
const (
	BannedReasonVPN        BannedReasonType = "VPN"
	BannedReasonMiner      BannedReasonType = "MINER"
	BannedReasonAggregator BannedReasonType = "AGGREGATOR"
	BannedReasonMirror     BannedReasonType = "MIRROR"
	BannedReasonDMCA       BannedReasonType = "DMCA"
	BannedReasonIllegal    BannedReasonType = "ILLEGAL"
)

// User is the simplest model of user, which is used in most queries.
type User struct {
	CreatedAt *time.Time `json:"createdAt" graphql:"createdAt"`
	BannedAt  *time.Time `json:"bannedAt" graphql:"bannedAt"`
	// BannedReason    *BannedReasonType  `json:"bannedReason" graphql:"bannedReason"`
	// EmailPreference map[string]bool    `json:"emailPreference" graphql:"emailPreference"`
	AgreedAt *time.Time `json:"agreedAt" graphql:"agreedAt"`
	// DiscordID is the user's Discord ID.
	DiscordID *string `json:"discordID" graphql:"discordID"`
	ID        string  `json:"_id" graphql:"_id"`
	Name      string  `json:"name" graphql:"name"`
	Email     string  `json:"email" graphql:"email"`
	Username  string  `json:"username" graphql:"username"`
	Language  string  `json:"language" graphql:"language"`
	AvatarURL string  `json:"avatarUrl" graphql:"avatarURL"`
	GitHubID  int64   `json:"githubID" graphql:"githubID"`
}

func (u *User) Header() []string {
	return []string{"ID", "Name", "Username", "Email", "Language", "RegisteredAt"}
}

func (u *User) Rows() [][]string {
	row := make([]string, 0, len(u.Header()))
	row = append(row, u.ID)
	row = append(row, u.Name)
	row = append(row, u.Username)
	row = append(row, u.Email)
	row = append(row, u.Language)
	row = append(row, util.ConvertTimeAgoString(*u.CreatedAt))

	return [][]string{row}
}

var _ Tabler = (*User)(nil)
