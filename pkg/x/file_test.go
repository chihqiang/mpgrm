package x

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMatchedFiles(t *testing.T) {
	// 创建临时测试目录
	tmpDir := filepath.Join(os.TempDir(), "mpgrm")
	log.Printf("temp %s", tmpDir)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("failed to create tmp dir: %v", err)
	}
	// 测试完成后删除目录
	//defer os.RemoveAll(tmpDir)

	files := []string{"a.txt", "b.txt", "c.log", "d.txt"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", f, err)
		}
	}

	tests := []struct {
		name     string
		patterns []string
		want     []string
	}{
		{
			"match txt files",
			[]string{filepath.Join(tmpDir, "*.txt")},
			[]string{
				filepath.Join(tmpDir, "a.txt"),
				filepath.Join(tmpDir, "b.txt"),
				filepath.Join(tmpDir, "d.txt"),
			},
		},
		{
			"match all files",
			[]string{filepath.Join(tmpDir, "*")},
			[]string{
				filepath.Join(tmpDir, "a.txt"),
				filepath.Join(tmpDir, "b.txt"),
				filepath.Join(tmpDir, "c.log"),
				filepath.Join(tmpDir, "d.txt"),
			},
		},
		{
			"empty pattern",
			[]string{""},
			[]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchedFiles(tt.patterns)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchedFiles(%v) = %v, want %v", tt.patterns, got, tt.want)
			}
		})
	}
}
