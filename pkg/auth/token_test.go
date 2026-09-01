package auth_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/zeabur/cli/pkg/auth"
)

func TestIsLegacyAPIKey(t *testing.T) {
	cases := map[string]bool{
		"sk-abcdefghijklmnopqrstuvwxyz012":       true,
		"zat_5f1a_abcdefghijklmnopqrstuvwxyz012": false,
		"eyJhbGciOiJIUzI1NiJ9.x.y":               false,
		"":                                       false,
	}
	for token, want := range cases {
		if got := auth.IsLegacyAPIKey(token); got != want {
			t.Errorf("IsLegacyAPIKey(%q) = %v, want %v", token, got, want)
		}
	}
}

func TestCallbackServerAcceptsTokenFields(t *testing.T) {
	cases := []struct {
		name string
		form url.Values
		want string
	}{
		{"access_token preferred", url.Values{"access_token": {"zat_1"}, "api_key": {"sk-1"}}, "zat_1"},
		{"api_key fallback", url.Values{"api_key": {"sk-1"}}, "sk-1"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server, err := auth.NewCallbackServer()
			if err != nil {
				t.Fatal(err)
			}
			go func() { _ = server.Serve() }()
			defer func() { _ = server.Close() }()

			c.form.Set("state", "abc")
			resp, err := http.Post(server.GetCallbackURL(), "application/x-www-form-urlencoded", strings.NewReader(c.form.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			_ = resp.Body.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			got, err := server.WaitForToken(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if got.Token != c.want || got.State != "abc" {
				t.Errorf("got %+v, want token %q state abc", got, c.want)
			}
		})
	}
}
