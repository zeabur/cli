package api

import (
	"github.com/zeabur/cli/pkg/model"
)

// V means graphql variables, it's a alias of map[string]interface{}
type V map[string]interface{}

type MapString map[string]string

func (id MapString) GetGraphQLType() string {
	return `Map`
}

// ObjectID is the alias ofskip, limit = normalizePagination(skip, limit) string, it's used to represent the ObjectID defined in GraphQL schema.
type ObjectID string

// GetGraphQLType returns the GraphQL type name of ObjectID.
func (id ObjectID) GetGraphQLType() string {
	return `ObjectID`
}

type ServiceTemplate string

func (t ServiceTemplate) GetGraphQLType() string {
	return `ServiceTemplate`
}

type GitProvider string

func (g GitProvider) GitProvider() string {
	return `GitProvider`
}

func normalizePagination(skip, limit int) (int, int) {
	if skip < 0 {
		skip = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 5
	}
	return skip, limit
}

type queryWithPagination[T any] func(skip, limit int) (*model.Connection[T], error)

// listAll is a helper function to list all items from a paginated query.
func listAll[T any](query queryWithPagination[T]) (items []*T, err error) {
	skip := 0
	perPage := 10

	next := true

	for next {
		itemCon, err := query(skip, perPage)
		if err != nil {
			return nil, err
		}
		for _, item := range itemCon.Edges {
			items = append(items, item.Node)
		}

		skip += perPage
		next = itemCon.PageInfo.HasNextPage
	}

	return items, nil
}
