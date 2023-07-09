package api

import (
	"context"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

const ZeaburGraphQLAPIEndpoint = "https://gateway.zeabur.dev/graphql"

type client struct {
	*graphql.Client
}

// New returns a new Zeabur API client.
func New(token string) Client {
	return &client{
		NewGraphQLClientWithToken(token),
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

var _ Client = (*client)(nil)
