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

func maskToken(token string) string {
	length := len(token)
	if length <= 8 {
		return strings.Repeat("*", length)
	}
	return token[:4] + "****" + token[length-4:]
}

// 实现 fmt.Stringer 接口
func (c Credential) String() string {
	// 脱敏 Token，只显示前 4 位，剩余部分用 * 替代
	maskedToken := maskToken(c.Token)
	return fmt.Sprintf("Username: %s, Token: %s, CloneURL: %s",
		c.Username, maskedToken, c.CloneURL)
}

// SetCloneByRepoName
// 当CloneURL为组织的地址时候才去调用补充完整的git地址进行clone
// SetCloneByRepoName 根据仓库名补充完整的 Git 仓库 CloneURL
//
// 当 CloneURL 是组织地址时（不是完整仓库地址），
// 会将 repoName 拼接到 CloneURL 后，生成完整仓库地址，并保证以 .git 结尾。
// 如果 CloneURL 已经是完整仓库地址（以 .git 结尾），不会进行修改。
//
// 参数:
//
//	repoName - 仓库名称，可以带或不带 .git，函数会自动处理
//
// 返回值:
//
//	error - 当 CloneURL 为空时返回错误，否则返回 nil
func (c *Credential) SetCloneByRepoName(repoName string) error {
	// 如果 CloneURL 为空，则无法拼接仓库地址，返回错误
	if c.CloneURL == "" {
		return fmt.Errorf("CloneURL is empty, cannot set repository")
	}
	// 如果 CloneURL 已经是完整仓库地址（以 .git 结尾），说明不需要补充仓库名，直接返回
	if strings.HasSuffix(c.CloneURL, ".git") {
		return nil
	}
	// 清理 repoName：去掉前后空格
	repoName = strings.TrimSpace(repoName)
	// 去掉 repoName 末尾的 .git（避免重复）
	repoName = strings.TrimSuffix(repoName, ".git")
	// 去掉 repoName 末尾的斜杠 /
	repoName = strings.TrimSuffix(repoName, "/")
	// 清理 CloneURL：去掉末尾的斜杠 /（避免拼接时出现双斜杠）
	cloneURL := strings.TrimSuffix(c.CloneURL, "/")
	// 拼接完整仓库地址，并保证以 .git 结尾
	c.CloneURL = fmt.Sprintf("%s/%s.git", cloneURL, repoName)
	return nil
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
