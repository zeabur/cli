package auth

import "strings"

const (
	// AccessTokenPrefix marks a scoped personal access token issued by createAccessToken.
	AccessTokenPrefix = "zat_"
	// LegacyAPIKeyPrefix marks an unscoped API key issued by the deprecated createApiKey.
	LegacyAPIKeyPrefix = "sk-"
)

const LegacyAPIKeyDeprecationMessage = "You are authenticated with a legacy API key, which is deprecated. " +
	"Run `zeabur auth login` to migrate to an access token, " +
	"or create one at " + ZeaburAccessTokenSettingsURL + " for use with ZEABUR_TOKEN."

func IsLegacyAPIKey(token string) bool {
	return strings.HasPrefix(token, LegacyAPIKeyPrefix)
}
