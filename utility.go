package fastci

import (
	"regexp"
	"strings"
)

var (
	regexpNonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

func sanitizeEnvKey(name string) string {
	name = strings.ToUpper(name)
	name = regexpNonAlphanumeric.ReplaceAllString(name, "_")
	return name
}
