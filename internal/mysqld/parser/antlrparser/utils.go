package antlrparser

import "strings"

func CaseInsensitiveEqual(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}
