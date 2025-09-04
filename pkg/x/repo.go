package x

import (
	"fmt"
	"net/url"
	"strings"
)

// RepoParseFullName extracts the repository owner and repository name from FullName.
// FullName format is usually "owner/repo", but may contain multiple path segments.
// Rules:
// - owner: the first segment
// - repo: the remaining segments joined by "/"
// Examples:
//
//	"openai/chatgpt"            -> owner="openai", repo="chatgpt"
//	"google/cloud/storage"      -> owner="google", repo="cloud/storage"
//	"microsoft/azure/devops"    -> owner="microsoft", repo="azure/devops"
func RepoParseFullName(fullName string) (owner string, repo string, err error) {
	if fullName == "" {
		return "", "", fmt.Errorf("FullName is empty")
	}
	fullNameParts := strings.Split(fullName, "/")
	if len(fullNameParts) < 2 {
		return "", "", fmt.Errorf("invalid FullName format: %s", fullName)
	}
	owner = fullNameParts[0]
	repo = strings.Join(fullNameParts[1:], "/") // join remaining parts as repo
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

func RepoFullNameOrgName(path string) string {
	// 去掉开头和结尾的 /
	path = strings.Trim(path, "/")
	if path == "" {
		return ""
	}
	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		// 只有一段，直接返回
		return parts[0]
	}
	// 判断最后一段是否是仓库名（通常带 .git 或者当作 repo）
	last := parts[len(parts)-1]
	if strings.HasSuffix(last, ".git") || len(parts) > 1 {
		// 去掉最后一段
		return strings.Join(parts[:len(parts)-1], "/")
	}
	// 否则返回完整路径
	return path
}
