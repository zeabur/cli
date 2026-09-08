package value

const visiblePrefixLength = 3

// Mask hides a variable value for display while preserving the short prefix
// used by the existing interactive variable-update flow.
func Mask(raw string) string {
	runes := []rune(raw)
	if len(runes) <= visiblePrefixLength {
		return "***"
	}
	return string(runes[:visiblePrefixLength]) + "***"
}
