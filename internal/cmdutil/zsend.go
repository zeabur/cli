package cmdutil

import (
	"fmt"
	"os"
	"strings"
)

// ZSendAPIKeyEnv is the environment variable read when --api-key is not set.
const ZSendAPIKeyEnv = "ZSEND_API_KEY"

// ZSendAPIKeyFlagUsage is the shared help text of the --api-key flag on
// Z-Send REST commands.
const ZSendAPIKeyFlagUsage = "Z-Send API key (or set ZSEND_API_KEY); defaults to the logged-in CLI token"

// ResolveZSendAPIKey returns the credential for the Z-Send REST API.
//
// Resolution order:
//  1. the --api-key flag value
//  2. the ZSEND_API_KEY environment variable
//  3. the logged-in CLI token, which the Z-Send REST API also accepts
//
// An explicitly supplied key (1 or 2) must be a Z-Send key with the zs_
// prefix; the CLI token is used as-is.
func (f *Factory) ResolveZSendAPIKey(flagValue string) (string, error) {
	apiKey := flagValue
	if apiKey == "" {
		apiKey = os.Getenv(ZSendAPIKeyEnv)
	}
	if apiKey != "" {
		if !strings.HasPrefix(apiKey, "zs_") {
			return "", fmt.Errorf("invalid API key format: must start with zs_")
		}
		return apiKey, nil
	}

	if f.Config != nil {
		if token := f.Config.GetTokenString(); token != "" {
			return token, nil
		}
	}

	return "", fmt.Errorf("Z-Send API key is required: pass --api-key, set ZSEND_API_KEY, or run `zeabur auth login`")
}
