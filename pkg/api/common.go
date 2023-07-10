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
