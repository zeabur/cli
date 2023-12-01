package auth

import "github.com/zeabur/cli/pkg/constant"

// Zeabur OAuth constants
const (
	ZeaburOAuthServerURL    = constant.ZeaburServerURL + "/oauth"
	ZeaburOAuthAuthorizeURL = ZeaburOAuthServerURL + "/authorize"
	ZeaburOAuthTokenURL     = ZeaburOAuthServerURL + "/token"

	ZeaburOAuthCLIClientID     = "c92dcff3-92a1-4d39-827e-a3eb952d5e0b"
	ZeaburOAuthCLIClientSecret = "tI2Q32jf"

	OAuthLocalServerCallbackURL = "http://localhost/callback"
)
