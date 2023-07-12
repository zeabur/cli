package api

import (
	"fmt"
)

// V means graphql variables, it's a alias of map[string]interface{}
type V map[string]interface{}

type ObjectID string

func (id ObjectID) GetGraphQLType() string {
	return fmt.Sprintf(`ObjectID`)
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
