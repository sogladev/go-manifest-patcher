package updater

import (
	"regexp"
	"strings"
)

var suffixRegex = regexp.MustCompile(`-\w+$`)

func MatchSpecialEdition(tagName string) bool {
	return !suffixRegex.MatchString(tagName)
}

func CompareVersions(v1, v2 string) int {
	return strings.Compare(v1, v2)
}
