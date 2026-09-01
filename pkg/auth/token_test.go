package auth

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestIsLegacyAPIKey(t *testing.T) {
	cases := map[string]bool{
		"sk-abcdefghijklmnopqrstuvwxyz012":       true,
		"zat_5f1a_abcdefghijklmnopqrstuvwxyz012": false,
		"eyJhbGciOiJIUzI1NiJ9.x.y":               false,
		"":                                       false,
	}
	for token, want := range cases {
		if got := IsLegacyAPIKey(token); got != want {
			t.Errorf("IsLegacyAPIKey(%q) = %v, want %v", token, got, want)
		}
	}
}

func TestTokenFromForm(t *testing.T) {
	cases := []struct {
		name string
		form url.Values
		want string
	}{
		{"access_token preferred", url.Values{"access_token": {"zat_1"}, "api_key": {"sk-1"}}, "zat_1"},
		{"api_key fallback", url.Values{"api_key": {"sk-1"}}, "sk-1"},
		{"empty", url.Values{}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodPost, "/callback", strings.NewReader(c.form.Encode()))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if got := tokenFromForm(r); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}
