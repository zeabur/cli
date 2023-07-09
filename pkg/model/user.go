package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BannedReasonType string

const (
	BannedReasonVPN        BannedReasonType = "VPN"
	BannedReasonMiner      BannedReasonType = "MINER"
	BannedReasonAggregator BannedReasonType = "AGGREGATOR"
	BannedReasonMirror     BannedReasonType = "MIRROR"
	BannedReasonDMCA       BannedReasonType = "DMCA"
	BannedReasonIllegal    BannedReasonType = "ILLEGAL"
)

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id" graphql:"_id"`
	Name      string             `json:"name" bson:"name" graphql:"name"`
	Email     string             `json:"email" bson:"email" graphql:"email"`
	Username  string             `json:"username" bson:"username" graphql:"username"`
	Language  string             `json:"language" bson:"language" graphql:"language"`
	GitHubID  int64              `json:"githubID" bson:"githubID" graphql:"githubID"`
	AvatarURL string             `json:"avatarUrl" bson:"avatarUrl" graphql:"avatarURL"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt" graphql:"createdAt"`
	BannedAt  *time.Time         `json:"bannedAt" bson:"bannedAt" graphql:"bannedAt"`
	//BannedReason    *BannedReasonType  `json:"bannedReason" bson:"bannedReason" graphql:"bannedReason"`
	//EmailPreference map[string]bool    `json:"emailPreference" bson:"emailPreference" graphql:"emailPreference"`
	AgreedAt        *time.Time `json:"agreedAt" bson:"agreedAt" graphql:"agreedAt"`
	LastCheckedInAt *time.Time `json:"lastCheckedInAt" bson:"lastCheckedInAt" graphql:"lastCheckedInAt"`
	// DiscordID is the user's Discord ID.
	DiscordID *string `json:"discordID" bson:"discordID" graphql:"discordID"`
}
