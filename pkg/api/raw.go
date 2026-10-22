package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// GraphQLRaw sends a pre-encoded GraphQL request body ({"query": ..., "variables": ...})
// to the Zeabur GraphQL API with the client's credentials and returns the
// response body verbatim, so callers can inspect `data` and `errors` as-is.
//
// A non-2xx HTTP status is returned as an error together with the body.
func (c *client) GraphQLRaw(ctx context.Context, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ZeaburGraphQLAPIEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBody, fmt.Errorf("graphql request failed: HTTP %d", resp.StatusCode)
	}

	return respBody, nil
}
