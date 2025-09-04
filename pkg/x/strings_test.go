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
