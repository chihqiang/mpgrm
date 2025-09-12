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

// HideSensitive 遮蔽敏感信息，只显示前后 l 个字符，中间用 **** 代替
func HideSensitive(str string, l int) string {
	length := len(str)
	if length == 0 {
		return ""
	}
	// 如果长度太短或 l*2 大于等于长度，全部用 *
	if length <= 6 || l*2 >= length {
		return strings.Repeat("*", length)
	}
	// 返回前 l 个字符 + **** + 后 l 个字符
	return str[:l] + "****" + str[length-l:]
}
