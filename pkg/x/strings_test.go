package x

import (
	"reflect"
	"testing"
)

func TestStringSplit(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		sep      string
		expected []string
	}{
		{
			name:     "normal split without whitespace",
			input:    "apple,banana,orange",
			sep:      ",",
			expected: []string{"apple", "banana", "orange"},
		},
		{
			name:     "elements with leading/trailing whitespace",
			input:    "  hello  ,  world  ,  go  ",
			sep:      ",",
			expected: []string{"hello", "world", "go"},
		},
		{
			name:     "consecutive separators",
			input:    "a,,b,,c",
			sep:      ",",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separators at start and end",
			input:    ",test,case,",
			sep:      ",",
			expected: []string{"test", "case"},
		},
		{
			name:     "no separators present",
			input:    "hello world",
			sep:      ",",
			expected: []string{"hello world"},
		},
		{
			name:     "space as separator",
			input:    "  a   b  c   ",
			sep:      " ",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "multi-character separator",
			input:    "x=1&&y=2&&z=3",
			sep:      "&&",
			expected: []string{"x=1", "y=2", "z=3"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StringSplit(tc.input, tc.sep)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("test failed: expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestStringSplitUniq(t *testing.T) {
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
			got := StringSplitUniq(tt.ss, tt.sep)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSplitUniq(%v, %q) = %v, want %v", tt.ss, tt.sep, got, tt.want)
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
