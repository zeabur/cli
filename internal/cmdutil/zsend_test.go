package cmdutil_test

import (
	"testing"

	"github.com/zeabur/cli/internal/cmdutil"
)

// tokenConfig is a stubConfig that reports a logged-in CLI token.
type tokenConfig struct {
	stubConfig
	token string
}

func (c tokenConfig) GetTokenString() string { return c.token }

func TestResolveZSendAPIKey(t *testing.T) {
	const cliToken = "cli-login-token"

	tests := []struct {
		name    string
		flag    string
		env     string
		token   string
		want    string
		wantErr string
	}{
		{name: "flag wins", flag: "zs_flag", env: "zs_env", token: cliToken, want: "zs_flag"},
		{name: "env when flag unset", env: "zs_env", token: cliToken, want: "zs_env"},
		{name: "falls back to CLI token", token: cliToken, want: cliToken},
		{name: "invalid flag key", flag: "sk_flag", token: cliToken, wantErr: "invalid API key format: must start with zs_"},
		{name: "invalid env key", env: "sk_env", token: cliToken, wantErr: "invalid API key format: must start with zs_"},
		{name: "nothing configured", wantErr: "Z-Send API key is required: pass --api-key, set ZSEND_API_KEY, or run `zeabur auth login`"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(cmdutil.ZSendAPIKeyEnv, tt.env)
			f := &cmdutil.Factory{Config: tokenConfig{token: tt.token}}

			got, err := f.ResolveZSendAPIKey(tt.flag)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("err = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
