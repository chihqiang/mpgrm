package x

import (
	"net/url"
	"testing"
)

func TestRepoParseFullName(t *testing.T) {
	tests := []struct {
		fullName  string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		// 普通 owner/repo
		{"openai/chatgpt", "openai", "chatgpt", false},
		// 多层 repo
		{"google/cloud/storage", "google/cloud", "storage", false},
		{"microsoft/azure/devops", "microsoft/azure", "devops", false},
		// 带尾部斜杠（不推荐，但测试兼容性）
		{"owner/repo/", "", "", true},
		// owner 为空
		{"/repo", "", "", true}, // 可以根据实际需求决定是否报错
		// repo 为空
		{"owner/", "", "", true}, // 可以根据实际需求决定是否报错
		// 格式错误
		{"invalidformat", "", "", true},
		// 空字符串
		{"", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.fullName, func(t *testing.T) {
			owner, repo, err := RepoParseFullName(tt.fullName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RepoParseFullName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if owner != tt.wantOwner {
				t.Errorf("RepoParseFullName() owner = %v, want %v", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("RepoParseFullName() repo = %v, want %v", repo, tt.wantRepo)
			}
		})
	}
}
func TestRepoURLParseFullName(t *testing.T) {
	tests := []struct {
		rawURL   string
		expected string
		wantErr  bool
	}{
		// HTTPS URL 带 .git
		{
			rawURL:   "https://github.com/owner/repo.git",
			expected: "owner/repo",
			wantErr:  false,
		},
		// HTTPS URL 多层目录
		{
			rawURL:   "https://github.com/owner/repo/aaa/aaa.git",
			expected: "owner/repo/aaa/aaa",
			wantErr:  false,
		},
		// GitLab HTTPS
		{
			rawURL:   "https://gitlab.com/group/project.git",
			expected: "group/project",
			wantErr:  false,
		},
		// HTTPS URL 末尾 /
		{
			rawURL:   "https://example.com/owner/repo/",
			expected: "owner/repo",
			wantErr:  false,
		},
		// 空 URL 应该报错
		{
			rawURL:   "",
			expected: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.rawURL, func(t *testing.T) {
			repoURL, _ := url.Parse(tt.rawURL)
			fullName, err := RepoURLParseFullName(repoURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RepoURLParseFullNameFromString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if fullName != tt.expected {
				t.Fatalf("RepoURLParseFullNameFromString() = %v, want %v", fullName, tt.expected)
			}
		})
	}
}

func TestRepoURLParseOrgName(t *testing.T) {
	tests := []struct {
		rawURL  string
		want    string
		wantErr bool
	}{
		{"https://github.com/org/repo.git", "org", false},
		{"https://github.com/org/org1/repo.git", "org/org1", false},
		{"https://github.com/org/org1/org2/repo.git", "org/org1/org2", false},
		{"https://github.com/org", "", false},
		{"https://github.com/org/ee", "org", false},
		{"https://github.com", "", false},
		{"https://github.com/", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.rawURL, func(t *testing.T) {
			u, _ := url.Parse(tt.rawURL)
			got, err := RepoURLParseOrgName(u)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepoURLParseOrgName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RepoURLParseOrgName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoFullNameOrgName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/org/repo.git", "org"},
		{"/org/org1/repo.git", "org/org1"},
		{"/org/org1/repo", "org/org1"},
		{"/org/org1/repo/", "org/org1/repo"},
		{"/org", ""},
		{"/", ""},
		{"", ""},
		{"org/repo.git", "org"},
		{"org/org1/repo", "org/org1"},
		{"org/org1/repo/", "org/org1/repo"},
		{"single", ""},
	}

	for _, tt := range tests {
		got := RepoFullNameOrgName(tt.input)
		if got != tt.expected {
			t.Errorf("RepoFullNameOrgName(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}
