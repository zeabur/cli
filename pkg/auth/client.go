// Package auth implements the authentication flow for the CLI
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/cli/browser"
)

type Client interface {
	GenerateToken(ctx context.Context) (string, error)
}

type ImplicitFlowClient struct {
	Endpoint url.URL

	callbackServer *CallbackServer
}

func NewImplicitFlowClient(callbackServer *CallbackServer) *ImplicitFlowClient {
	endpointURL, err := url.Parse(ZeaburApiKeyConfirmEndpoint)
	if err != nil {
		panic(fmt.Sprintf("failed to parse endpoint URL (internal error): %v", err))
	}

	return &ImplicitFlowClient{
		Endpoint:       *endpointURL,
		callbackServer: callbackServer,
	}
}

func (c *ImplicitFlowClient) GenerateToken(ctx context.Context) (token string, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Generate a state that is only used for the callback
	state, err := randomString(40)
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	go func() {
		_ = c.callbackServer.Serve()
	}()
	defer func() {
		_ = c.callbackServer.Close()
	}()

	// Construct the confirmation URL
	endpoint := c.Endpoint

	query := endpoint.Query()
	query.Add("client_name", "zeabur-cli")
	query.Add("state", state)
	query.Add("callback_url", c.callbackServer.GetCallbackURL())

	endpoint.RawQuery = query.Encode()

	// Open the browser
	if err := browser.OpenURL(endpoint.String()); err != nil {
		return "", fmt.Errorf("failed to open browser (url=%s): %w", endpoint.String(), err)
	}

	// Wait for the token
	tokenResponse, err := c.callbackServer.WaitForToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to wait for token: %w", err)
	}

	if tokenResponse.State != state {
		return "", fmt.Errorf("state mismatch: expected=%s, got=%s", state, tokenResponse.State)
	}

	return tokenResponse.Token, nil
}

func randomString(length int) (string, error) {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
