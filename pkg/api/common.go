package api

// V means graphql variables, it's a alias of map[string]interface{}
type V map[string]interface{}

// ObjectID is the alias of string, it's used to represent the ObjectID defined in GraphQL schema.
type ObjectID string

// GetGraphQLType returns the GraphQL type name of ObjectID.
func (id ObjectID) GetGraphQLType() string {
	return `ObjectID`
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
