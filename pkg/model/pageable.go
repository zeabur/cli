package model

type Connection[T any] struct {
	PageInfo *PageInfo  `json:"pageInfo" graphql:"pageInfo"`
	Edges    []*Edge[T] `json:"edges" graphql:"edges"`
}

type Edge[T any] struct {
	Node   *T     `json:"node" graphql:"node"`
	Cursor string `json:"cursor" graphql:"cursor"`
}

// PageInfo is the pagination information
type PageInfo struct {
	StartCursor     string `json:"startCursor" graphql:"startCursor"`
	EndCursor       string `json:"endCursor" graphql:"endCursor"`
	TotalCount      int    `json:"totalCount" graphql:"totalCount"`
	HasNextPage     bool   `json:"hasNextPage" graphql:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage" graphql:"hasPreviousPage"`
}
