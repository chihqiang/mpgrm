package platforms

import (
	"context"
	"fmt"
	"github.com/chihqiang/mpgrm/pkg/httpx"
	"github.com/chihqiang/mpgrm/pkg/logx"
	"github.com/chihqiang/mpgrm/pkg/x"
	"golang.org/x/sync/errgroup"
	"path"
	"sync"
)

type ReleaseInfo struct {
	ID          int64  // Release 的唯一 ID（由平台生成），用于获取或删除
	TagName     string // 关联的标签名，例如 "v1.0.0"
	Title       string // Release 的标题名称，例如 "Initial Release"
	Description string // Release 的详细描述内容（通常用于更新日志等），例如 "Added login feature, fixed bugs"

	FullName string // 仓库全名，通常是 "组织名/仓库名"，例如: "my-org/my-repo"

	Assets []*AssetInfo // 附件列表，表示该 Release 中包含的所有资源文件（如安装包、构建产物等）
}

// AssetInfo 表示一个发布版本（Release）中附带的单个附件信息。
// 常用于获取、上传或删除 Release 中的二进制包、安装文件等资源。
type AssetInfo struct {
	ID   int64  // 附件的唯一 ID（由平台生成），用于删除或引用
	Name string // 附件的名称（如文件名），如 "app-v1.0.0-linux-amd64.tar.gz"
	URL  string // 附件的下载 URL（通常是公开链接或 API 地址）
}

func (ri *ReleaseInfo) Init() {
	if ri.Title == "" {
		ri.Title = ri.TagName
	}
	if ri.Description == "" {
		ri.Description = ri.TagName
	}
}

func (ri *ReleaseInfo) Download(workspace string) ([]string, error) {
	var (
		localFileNames []string
		mu             sync.Mutex
	)
	g, ctx := errgroup.WithContext(context.Background())
	for _, asset := range ri.Assets {
		asset := asset
		g.Go(func() error {
			mu.Lock()
			localFile := path.Join(workspace, ri.TagName, asset.Name)
			logx.Debug("Starting download: %s", asset.URL)
			if err := httpx.Download(ctx, asset.URL, localFile); err != nil {
				return fmt.Errorf("failed to download %s: %w", asset.URL, err)
			}
			localFileNames = append(localFileNames, localFile)
			logx.Debug("Download completed: %s", localFile)
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return localFileNames, nil
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
func (ri *ReleaseInfo) GetOwnerRepo() (owner string, repo string, err error) {
	return x.RepoParseFullName(ri.FullName)
}
