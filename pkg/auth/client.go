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
	WebAppClient struct {
		ClientID               string
		ClientSecret           string
		RedirectURIWithoutPort string
		RedirectURIWithPort    string // we will fill port after local server is started
		Scopes                 []string

		AuthorizeURL string
		TokenURL     string

		config oauth2.Config

		httpClient *http.Client
	}

	Options struct {
		ClientID               string
		ClientSecret           string
		RedirectURIWithoutPort string
		Scopes                 []string

		AuthorizeURL string
		TokenURL     string
		HttpClient   *http.Client
	}
)

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

func NewWebAppClient(opts Options) *WebAppClient {
	c := &WebAppClient{
		ClientID:               opts.ClientID,
		ClientSecret:           opts.ClientSecret,
		RedirectURIWithoutPort: opts.RedirectURIWithoutPort,
		Scopes:                 opts.Scopes,

		AuthorizeURL: opts.AuthorizeURL,
		TokenURL:     opts.TokenURL,
	}

	if opts.HttpClient != nil {
		c.httpClient = opts.HttpClient
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

func (c *WebAppClient) RefreshToken(old *oauth2.Token) (new *oauth2.Token, err error) {
	new, err = c.config.TokenSource(context.Background(), old).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return new, nil
}

var _ Client = (*WebAppClient)(nil)
