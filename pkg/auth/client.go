// Package auth implements the authentication flow for the CLI
package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cli/browser"
	"golang.org/x/oauth2"

	"github.com/zeabur/cli/pkg/webapp"
)

type (
	// WebAppClient is a client for OAuth2 authorization code flow.
	WebAppClient struct {
		httpClient *http.Client

		ClientID               string
		ClientSecret           string
		RedirectURIWithoutPort string
		RedirectURIWithPort    string // we will fill port after local server is started
		AuthorizeURL           string
		TokenURL               string
		Scopes                 []string

		config oauth2.Config
	}

	// Options is the options for WebAppClient.
	Options struct {
		ClientID               string
		ClientSecret           string
		RedirectURIWithoutPort string

		AuthorizeURL string
		TokenURL     string
		HTTPClient   *http.Client
		Scopes       []string
	}
)

// NewZeaburWebAppOAuthClient creates a new WebAppClient for Zeabur OAuth server.
func NewZeaburWebAppOAuthClient() *WebAppClient {
	opts := Options{
		ClientID:               ZeaburOAuthCLIClientID,
		ClientSecret:           ZeaburOAuthCLIClientSecret,
		RedirectURIWithoutPort: OAuthLocalServerCallbackURL,
		Scopes:                 []string{"all"},

		AuthorizeURL: ZeaburOAuthAuthorizeURL,
		TokenURL:     ZeaburOAuthTokenURL,
	}

	return NewWebAppClient(opts)
}

// NewWebAppClient creates a new WebAppClient.
func NewWebAppClient(opts Options) *WebAppClient {
	c := &WebAppClient{
		ClientID:               opts.ClientID,
		ClientSecret:           opts.ClientSecret,
		RedirectURIWithoutPort: opts.RedirectURIWithoutPort,
		Scopes:                 opts.Scopes,

		AuthorizeURL: opts.AuthorizeURL,
		TokenURL:     opts.TokenURL,
	}

	if opts.HTTPClient != nil {
		c.httpClient = opts.HTTPClient
	} else {
		c.httpClient = http.DefaultClient
	}

	c.config = oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Scopes:       c.Scopes,
		RedirectURL:  "", // we will fill port after local server is started
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AuthorizeURL,
			TokenURL: c.TokenURL,
		},
	}

	return c
}

// Login helps the user to login to the OAuth server with a web browser.
func (c *WebAppClient) Login() (token *oauth2.Token, err error) {
	flow, err := webapp.InitFlow()
	if err != nil {
		panic(err)
	}

	redirectURIWithPort, err := flow.RedirectURIWithPort(c.RedirectURIWithoutPort)
	if err != nil {
		return nil, fmt.Errorf("failed to get redirect URI with port: %w", err)
	}

	c.config.RedirectURL = redirectURIWithPort

	browserURL, err := flow.BrowserURL(c.AuthorizeURL, c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to construct OAuth URL: %w", err)
	}

	// A localhost server on a random available port will receive the web redirect.
	go func() {
		_ = flow.StartServer(nil)
	}()

	//Note: the user's web browser must run on the same device as the running app.
	err = browser.OpenURL(browserURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open browser: %w", err)
	}

	accessToken, err := flow.Wait(context.TODO(), c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	return accessToken, nil
}

// RefreshToken refreshes the token.
func (c *WebAppClient) RefreshToken(old *oauth2.Token) (newToken *oauth2.Token, err error) {
	newToken, err = c.config.TokenSource(context.Background(), old).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}

var _ Client = (*WebAppClient)(nil)
