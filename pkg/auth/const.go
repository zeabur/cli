package auth

import "github.com/zeabur/cli/pkg/constant"

// Zeabur OAuth constants
const (
	ZeaburOAuthServerURL    = constant.ZeaburServerURL + "/oauth"
	ZeaburOAuthAuthorizeURL = ZeaburOAuthServerURL + "/authorize"
	ZeaburOAuthTokenURL     = ZeaburOAuthServerURL + "/token"

	ZeaburOAuthCLIClientID     = "64c7559d0e6da9ed35dee1ff"
	ZeaburOAuthCLIClientSecret = "630641"

	OAuthLocalServerCallbackURL = "http://localhost/callback"
)
