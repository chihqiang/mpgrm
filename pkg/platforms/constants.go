package platforms

import (
	"context"
	"wangzhiqiang/mpgrm/pkg/credential"
)

type IPlatform interface {
	credential.ICredential
	IPlatformRepo
	IPlatformReleases
}

// IPlatformRepo 定义了对代码托管平台仓库的操作接口。
// 包括列出仓库、获取单个仓库、创建和删除仓库等操作。
type IPlatformRepo interface {
	// ListOrgRepo 列出指定组织下的所有仓库。
	ListOrgRepo(ctx context.Context, orgName string) ([]*RepoInfo, error)

	// ListUserRepo 列出当前用户自己创建或拥有的仓库。
	ListUserRepo(ctx context.Context) ([]*RepoInfo, error)

	// GetRepoDetail 获取单个仓库的详细信息。
	GetRepoDetail(ctx context.Context, fullName string) (*RepoInfo, error)

	// CreateRepo 创建一个新的仓库。
	CreateRepo(ctx context.Context, repoInfo *RepoInfo) error

	// DeleteRepo 删除指定的仓库。
	DeleteRepo(ctx context.Context, repoInfo *RepoInfo) error
}

// IPlatformReleases 定义了平台发布版本相关的操作
type IPlatformReleases interface {
	// ListTags 列出指定仓库的所有标签
	ListTags(ctx context.Context, fullName string) ([]*TagInfo, error)

	// GetTagReleaseInfo 获取指定仓库下的某个标签信息
	GetTagReleaseInfo(ctx context.Context, fullName, tagName string) (*ReleaseInfo, error)

	// CreateRelease 在指定仓库下创建一个新的发布版本
	CreateRelease(ctx context.Context, fullName string, releaseInfo *ReleaseInfo) (newTagInfo *ReleaseInfo, er error)

	// DeleteReleaseAssets 删除指定发布版本下的一个或多个资源文件
	DeleteReleaseAssets(ctx context.Context, releaseInfo *ReleaseInfo, filenames []string) error

	// UploadReleaseAsset 上传一个或多个资源文件到指定的发布版本
	UploadReleaseAsset(ctx context.Context, releaseInfo *ReleaseInfo, filenames []string) error
}
