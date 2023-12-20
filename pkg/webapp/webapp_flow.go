// Package webapp implements the OAuth Web Application authorization flow for client applications by
// starting a server at localhost to receive the web redirect after the user has authorized the application.
package webapp

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

// Flow holds the state for the steps of OAuth Web Application flow.
type Flow struct {
	server        *localServer
	state         string
	codeChallenge string
}

// InitFlow creates a new Flow instance by detecting a locally available port number.
func InitFlow() (*Flow, error) {
	server, err := bindLocalServer()
	if err != nil {
		return nil, err
	}

	state, _ := randomString(20)
	codeChallenge, _ := randomString(43)

	return &Flow{
		server:        server,
		state:         state,
		codeChallenge: codeChallenge,
	}, nil
}

// RedirectURIWithPort completes the redirect url with the port number of the local server.
func (flow *Flow) RedirectURIWithPort(redirectURIWithoutPort string) (string, error) {
	ru, err := url.Parse(redirectURIWithoutPort)
	if err != nil {
		return "", fmt.Errorf("invalid redirect URI: %w", err)
	}

	ru.Host = fmt.Sprintf("%s:%d", ru.Hostname(), flow.server.Port())
	flow.server.CallbackPath = ru.Path

	return ru.String(), nil
}

// BrowserURL appends GET query parameters to baseURL and returns the url that the user should
// navigate to in their web browser.
func (flow *Flow) BrowserURL(baseURL string, config oauth2.Config) (string, error) {
	q := url.Values{}
	q.Set("client_id", config.ClientID)
	q.Set("redirect_uri", config.RedirectURL)
	q.Set("code_challenge", genCodeChallengeS256(flow.codeChallenge))
	q.Set("code_challenge_method", "S256")
	q.Set("scope", strings.Join(config.Scopes, ","))
	q.Set("state", flow.state)
	q.Set("response_type", "code")

	return fmt.Sprintf("%s?%s", baseURL, q.Encode()), nil
}

// StartServer starts the localhost server and blocks until it has received the web redirect. The
// writeSuccess function can be used to render a HTML page to the user upon completion.
func (flow *Flow) StartServer(writeSuccess func(io.Writer)) error {
	flow.server.WriteSuccessHTML = writeSuccess
	return flow.server.Serve()
}

// Wait blocks until the browser flow has completed and returns the access token.
func (flow *Flow) Wait(ctx context.Context, config oauth2.Config) (*oauth2.Token, error) {
	code, err := flow.server.WaitForCode(ctx)
	if err != nil {
		return nil, err
	}
	if code.State != flow.state {
		return nil, errors.New("state mismatch")
	}

	token, err := config.Exchange(context.Background(), code.Code, oauth2.SetAuthURLParam("code_verifier", flow.codeChallenge))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func randomString(length int) (string, error) {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
