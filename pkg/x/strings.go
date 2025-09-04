package x

import (
	"github.com/samber/lo"
	"strings"
)

// StringSplit splits a string by sep and trims each part, ignoring empty strings
func StringSplit(s, sep string) []string {
	var sp []string
	for _, p := range strings.Split(s, sep) {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			sp = append(sp, trimmed)
		}
	}
	return lo.Uniq(sp)
}

// StringSplits splits each string in ss by sep and flattens the result
func StringSplits(ss []string, sep string) []string {
	var sp []string
	for _, s := range ss {
		sp = append(sp, StringSplit(s, sep)...)
	}
	return lo.Uniq(sp)
}
