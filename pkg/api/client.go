// Package api provides a client for the Zeabur API.
package api

import (
	"context"
	"net/http"

	"github.com/zeabur/cli/pkg/constant"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

// ZeaburGraphQLAPIEndpoint is the endpoint for the Zeabur GraphQL API.
const (
	ZeaburGraphQLAPIEndpoint    = constant.ZeaburServerURL + "/graphql"
	WSSZeaburGraphQLAPIEndpoint = constant.WebsocketURL + "/graphql"
)

type client struct {
	*graphql.Client
	sub *graphql.SubscriptionClient

	// httpClient is the authenticated HTTP client shared with the GraphQL
	// client; used by GraphQLRaw to send caller-supplied request bodies.
	httpClient *http.Client
}

// New returns a new Zeabur API client.
func New(token string) Client {
	httpClient := newHTTPClientWithToken(token)
	return &client{
		Client:     graphql.NewClient(ZeaburGraphQLAPIEndpoint, httpClient),
		sub:        NewSubscriptionClient(token),
		httpClient: httpClient,
	}
}

// NewGraphQLClientWithToken returns a new GraphQL client with the given token.
func NewGraphQLClientWithToken(token string) *graphql.Client {
	return graphql.NewClient(ZeaburGraphQLAPIEndpoint, newHTTPClientWithToken(token))
}

func newHTTPClientWithToken(token string) *http.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	return oauth2.NewClient(context.Background(), src)
}

func NewSubscriptionClient(token string) *graphql.SubscriptionClient {
	return graphql.NewSubscriptionClient(WSSZeaburGraphQLAPIEndpoint).
		WithProtocol(graphql.GraphQLWS).
		WithConnectionParams(map[string]any{
			"authToken": token,
		})
}

var _ Client = (*client)(nil)
