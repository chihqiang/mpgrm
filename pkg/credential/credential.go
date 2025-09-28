package credential

import (
	"fmt"
	"github.com/chihqiang/mpgrm/pkg/x"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"net/url"
	"os"
	"path"
	"strings"
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
	Username string // 用于 HTTPS 认证的用户名，通常是 git 或你的账号名
	Password string // 用于 HTTPS 认证的密码（或者 Gitee 账号密码），用于拉取仓库
	Token    string // 用于 HTTPS 认证的访问令牌（Token），优先于密码
	CloneURL string // 仓库的克隆地址（HTTPS URL），用于 git clone
}

func (c *Credential) GetGitAuth() transport.AuthMethod {
	pwd := c.Password
	if pwd == "" {
		pwd = c.Token
	}
	if c.Username != "" && pwd != "" {
		return &http.BasicAuth{
			Username: c.Username,
			Password: pwd,
		}
	}
	return nil
}

// 实现 fmt.Stringer 接口
func (c Credential) String() string {
	return fmt.Sprintf(
		"Username: %s, Password: %s, Token: %s, CloneURL: %s",
		c.Username,
		x.HideSensitive(c.Password, 2),
		x.HideSensitive(c.Token, 3),
		c.CloneURL,
	)
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
