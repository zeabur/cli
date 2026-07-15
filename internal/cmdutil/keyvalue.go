package cmdutil

import (
	"fmt"
	"strings"
)

// ParseKeyValuePairs converts KEY=VALUE arguments into a map without
// interpreting characters in the value.
func ParseKeyValuePairs(values []string) (map[string]string, error) {
	pairs := make(map[string]string, len(values))

	for _, value := range values {
		key, val, found := strings.Cut(value, "=")
		if !found || key == "" {
			return nil, fmt.Errorf("invalid key value pair %q: expected KEY=VALUE format", value)
		}
		pairs[key] = val
	}

	return pairs, nil
}
