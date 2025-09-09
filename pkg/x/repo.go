package x

import (
	"fmt"
	"net/url"
	"strings"
)

// RepoParseFullName extracts the repository owner and repository name from FullName.
// Rules:
// - repo: the last segment, must not be empty
// - owner: all preceding segments joined by "/"
// - Trailing "/" is invalid
// - If FullName is invalid (empty, only one segment, last segment empty, or ends with "/"), return an error
// Examples:
//
//	"openai/chatgpt"        -> owner="openai", repo="chatgpt"
//	"google/cloud/storage"  -> owner="google/cloud", repo="storage"
//
// Invalid cases:
//
//	"owner/repo/" -> error
//	"/repo"       -> error
//	"owner/"      -> error
//	"invalidformat" -> error
//	""            -> error
func RepoParseFullName(fullName string) (owner string, repo string, err error) {
	if fullName == "" {
		return "", "", fmt.Errorf("FullName is empty")
	}

	// Trailing slash is invalid
	if strings.HasSuffix(fullName, "/") {
		return "", "", fmt.Errorf("invalid FullName, trailing slash not allowed: %s", fullName)
	}

	// Trim leading slashes
	fullName = strings.TrimLeft(fullName, "/")

	// Split path by "/"
	parts := strings.Split(fullName, "/")

	// Must have at least 2 segments
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid FullName format: %s", fullName)
	}

	// Last segment is repo
	repo = parts[len(parts)-1]
	if repo == "" {
		return "", "", fmt.Errorf("invalid FullName, repository name is empty: %s", fullName)
	}

	// Remove .git suffix if present
	repo = strings.TrimSuffix(repo, ".git")

	// Owner is all preceding segments joined by "/"
	owner = strings.Join(parts[:len(parts)-1], "/")
	if owner == "" {
		return "", "", fmt.Errorf("invalid FullName, owner is empty: %s", fullName)
	}

	return owner, repo, nil
}

// RepoURLParseFullName 从 repoURL 中解析仓库完整路径
// 返回格式: owner/repo[/子目录]
func RepoURLParseFullName(repoURL *url.URL) (string, error) {
	if repoURL == nil {
		return "", fmt.Errorf("repoURL is nil")
	}
	// 去掉开头的 /
	pathPart := strings.TrimPrefix(repoURL.Path, "/")
	// 去掉尾部的 .git
	pathPart = strings.TrimSuffix(pathPart, ".git")
	// 去掉尾部的 /
	pathPart = strings.TrimSuffix(pathPart, "/")
	return pathPart, nil
}

// RepoURLParseOrgName 从 repoURL 中解析组织/所有者名称
// 返回格式: owner[/子目录]，不包含最后的仓库名
func RepoURLParseOrgName(repoURL *url.URL) (string, error) {
	if repoURL == nil {
		return "", fmt.Errorf("repoURL is nil")
	}
	return RepoFullNameOrgName(repoURL.Path), nil
}

// RepoFullNameOrgName returns the organization path from a repository path by removing the repository name.
// Examples:
// 1. "/org/repo.git"       -> "org"           (.git suffix, remove last segment)
// 2. "/org/org1/repo.git"  -> "org/org1"      (.git suffix, remove last segment)
// 3. "/org/org1/repo"      -> "org/org1"      (no trailing /, remove last segment)
// 4. "/org/org1/repo/"     -> "org/org1/repo" (trailing /, keep last segment)
// 5. "/org"                -> ""              (single segment, return empty)
// 6. "/"                   -> ""              (empty path)
// 7. ""                    -> ""              (empty string)
func RepoFullNameOrgName(p string) string {
	if p == "" {
		return ""
	}
	// 拆分路径并去掉空片段（前导 / 会产生空字符串）
	parts := strings.Split(p, "/")
	cleanParts := []string{}
	for _, part := range parts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	if len(cleanParts) <= 1 {
		// 只有一段或空，返回空字符串
		return ""
	}

	last := cleanParts[len(cleanParts)-1]
	if strings.HasSuffix(last, ".git") {
		// 以 .git 结尾，去掉最后一段
		cleanParts = cleanParts[:len(cleanParts)-1]
	} else if !strings.HasSuffix(p, "/") {
		// 不以 / 结尾，去掉最后一段
		cleanParts = cleanParts[:len(cleanParts)-1]
	}
	return strings.Join(cleanParts, "/")
}
