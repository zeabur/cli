package auth

import "github.com/zeabur/cli/pkg/constant"

const (
	// ZeaburAccessTokenConfirmEndpoint is the dashboard page that mints an access token for the CLI.
	// The path keeps its historical "api-key" name so older CLI builds continue to resolve it.
	ZeaburAccessTokenConfirmEndpoint = constant.ZeaburDashURL + "/auth/api-key/confirm"

	ZeaburAccessTokenSettingsURL = constant.ZeaburDashURL + "/account/api-keys"
)
