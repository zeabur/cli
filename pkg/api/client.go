// Package api provides a client for the Zeabur API.
package api

import (
	"context"
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
}

// New returns a new Zeabur API client.
func New(token string) Client {
	return &client{
		NewGraphQLClientWithToken(token),
		NewSubscriptionClient(token),
	}
}

// NewGraphQLClientWithToken returns a new GraphQL client with the given token.
func NewGraphQLClientWithToken(token string) *graphql.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	return graphql.NewClient(ZeaburGraphQLAPIEndpoint, httpClient)
}

func NewSubscriptionClient(token string) *graphql.SubscriptionClient {
	return graphql.NewSubscriptionClient(ZeaburGraphQLAPIEndpoint).
		WithConnectionParams(map[string]any{
			"Authorization": "Bearer " + token,
		})
}

var _ Client = (*client)(nil)
