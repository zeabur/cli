package auth

import "github.com/zeabur/cli/pkg/constant"

// Zeabur OAuth constants
const (
	ZeaburOAuthServerURL    = constant.ZeaburServerURL + "/oauth"
	ZeaburOAuthAuthorizeURL = ZeaburOAuthServerURL + "/authorize"
	ZeaburOAuthTokenURL     = ZeaburOAuthServerURL + "/token"

	ZeaburOAuthCLIClientID     = "64880903673438f7c30bcc87"
	ZeaburOAuthCLIClientSecret = "999999"

	OAuthLocalServerCallbackURL = "http://localhost/callback"
)
