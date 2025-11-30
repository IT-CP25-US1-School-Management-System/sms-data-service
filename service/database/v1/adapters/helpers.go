package adapters

import "strings"

// joinStrings joins string slices with a separator
func joinStrings(parts []string, sep string) string {
	return strings.Join(parts, sep)
}

// trimQuotes removes surrounding single quotes from a string
func trimQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}
	return s
}
