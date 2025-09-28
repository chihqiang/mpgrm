package platforms

import (
	"fmt"
	"github.com/chihqiang/mpgrm/pkg/x"
)

// RepoInfo 表示代码托管平台上的仓库信息
type RepoInfo struct {
	ID          int64  // 仓库唯一 ID，例如: 123456
	Name        string // 仓库名称（短名），例如: "my-repo"
	FullName    string // 仓库全名，通常是 "组织名/仓库名"，例如: "my-org/my-repo"
	Description string // 仓库描述信息，例如: "这是一个示例仓库"
	Homepage    string // 仓库主页链接，例如: "https://example.com"
	IsPrivate   bool   // 是否为私有仓库，例如: true
	CloneURL    string // 克隆仓库的 URL（HTTPS 或 SSH），例如: "https://github.com/my-org/my-repo.git"
}

// GetOwnerRepo extracts the repository owner and repository name from FullName.
// FullName format is usually "owner/repo", but may contain multiple path segments.
// Rules:
// - owner: the first segment
// - repo: the remaining segments joined by "/"
// Examples:
//
//	"openai/chatgpt"            -> owner="openai", repo="chatgpt"
//	"google/cloud/storage"      -> owner="google", repo="cloud/storage"
//	"microsoft/azure/devops"    -> owner="microsoft", repo="azure/devops"
func (ri *RepoInfo) GetOwnerRepo() (owner string, repo string, err error) {
	return x.RepoParseFullName(ri.FullName)
}

func (ri *RepoInfo) GetOrgName() (string, error) {
	if ri.FullName == "" {
		return "", fmt.Errorf("RepoInfo FullName is null")
	}
	return x.RepoFullNameOrgName(ri.FullName), nil
}
