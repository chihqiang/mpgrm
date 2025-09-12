package x

import (
	"reflect"
	"testing"
)

func TestStringSplit(t *testing.T) {
	tests := []struct {
		name string
		s    string
		sep  string
		want []string
	}{
		{"normal split", "a,b,c", ",", []string{"a", "b", "c"}},
		{"trim spaces", " a , b ,c ", ",", []string{"a", "b", "c"}},
		{"empty parts", "a,,b,", ",", []string{"a", "b"}},
		{"all empty", ", , ", ",", []string{}},
		{"duplicates", "a,b,a", ",", []string{"a", "b"}},
		{"no sep", "abc", ",", []string{"abc"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringSplit(tt.s, tt.sep)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSplit(%q, %q) = %v, want %v", tt.s, tt.sep, got, tt.want)
			}
		})
	}
}

func TestStringSplits(t *testing.T) {
	tests := []struct {
		name string
		ss   []string
		sep  string
		want []string
	}{
		{"multiple strings", []string{"a,b", "c,d"}, ",", []string{"a", "b", "c", "d"}},
		{"with spaces", []string{" a , b ", " c "}, ",", []string{"a", "b", "c"}},
		{"empty strings", []string{"", " "}, ",", []string{}},
		{"duplicates across strings", []string{"a,b", "b,c"}, ",", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringSplits(tt.ss, tt.sep)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSplits(%v, %q) = %v, want %v", tt.ss, tt.sep, got, tt.want)
			}
		})
	}
}

func TestHideSensitive(t *testing.T) {
	tests := []struct {
		input   string
		visible int
		want    string
	}{
		{"", 2, ""},                       // 空字符串
		{"abc", 1, "***"},                 // 短字符串，全遮蔽
		{"abcdef", 2, "******"},           // 长度等于阈值，全遮蔽
		{"abcdefgh", 2, "ab****gh"},       // 正常遮蔽，显示前后2个字符
		{"abcdefghijkl", 3, "abc****jkl"}, // 正常遮蔽，显示前后3个字符
		{"abcdefgh", 5, "********"},       // l*2 > length，全遮蔽
	}

	for _, tt := range tests {
		got := HideSensitive(tt.input, tt.visible)
		if got != tt.want {
			t.Errorf("HideSensitive(%q, %d) = %q; want %q", tt.input, tt.visible, got, tt.want)
		}
	}
}
