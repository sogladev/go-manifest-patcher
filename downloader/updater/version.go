package updater

import (
	"regexp"
	"strings"
)

var suffixRegex = regexp.MustCompile(`-epoch`)

func MatchSpecialEdition(tagName string) bool {
	return suffixRegex.MatchString(tagName)
}

func CompareVersions(v1, v2 string) int {
	v1 = strings.TrimSuffix(v1, "-epoch")
	v2 = strings.TrimSuffix(v2, "-epoch")
	return strings.Compare(v1, v2)
}
