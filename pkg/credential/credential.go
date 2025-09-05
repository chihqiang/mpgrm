package credential

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"wangzhiqiang/mpgrm/pkg/x"
)

// WorkspaceCategory 表示 workspace 下的分类/命名空间层级
type WorkspaceCategory string

const (
	WorkspaceCategoryGit      WorkspaceCategory = "git"      // git 仓库
	WorkspaceCategoryReleases WorkspaceCategory = "releases" // release 文件
)

type ICredential interface {
	WithCredential(credential *Credential) error
}

type IFormToCredential interface {
	WithForm(credential *Credential)
	WithTarget(credential *Credential)
}

type Credential struct {
	Username string // 用户名或 git
	Token    string // 认证令牌（Token）
	CloneURL string // clone URL
}

// GetFullName 解析 CloneURL 并返回仓库完整名称，例如 "owner/repo"
func (c *Credential) GetFullName() (string, error) {
	if c.CloneURL == "" {
		return "", nil // 如果 CloneURL 为空，直接返回空字符串
	}
	parse, err := url.Parse(c.CloneURL)
	if err != nil {
		return "", err
	}
	return x.RepoURLParseFullName(parse)
}

func (c *Credential) GetCategoryNamWorkspace(category WorkspaceCategory, workspace string) (string, error) {
	if c.CloneURL == "" {
		return "", fmt.Errorf("CloneURL is empty")
	}
	// 将仓库地址转成安全目录名，比如去掉协议、替换斜杠为下划线
	u, err := url.Parse(c.CloneURL)
	if err != nil {
		return "", fmt.Errorf("invalid CloneURL: %w", err)
	}
	dirName := u.Host + u.Path
	dirName = strings.TrimSuffix(dirName, ".git") // 去掉 .git 后缀
	// 拼接 workspace
	localPath := path.Join(workspace, string(category), dirName)
	// 创建目录（如果不存在）
	if err := os.MkdirAll(localPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create workspace: %w", err)
	}
	return localPath, nil
}
