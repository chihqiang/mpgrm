package github

import (
	"context"
	"fmt"
	"github.com/chihqiang/mpgrm/pkg/httpx"
	"github.com/chihqiang/mpgrm/pkg/platforms"
	"github.com/chihqiang/mpgrm/pkg/x"
	"github.com/google/go-github/v73/github"
	"os"
	"path"
	"path/filepath"
)

func (p *Platform) ListTags(ctx context.Context, fullName string) ([]*platforms.TagInfo, error) {
	owner, repo, err := x.RepoParseFullName(fullName)
	if err != nil {
		return nil, err
	}
	client := p.GetClient(ctx)
	var (
		tagInfos []*platforms.TagInfo
	)
	err = httpx.Paginate[*github.RepositoryTag](func(page int) ([]*github.RepositoryTag, error) {
		tags, _, err := client.Repositories.ListTags(ctx, owner, repo,
			&github.ListOptions{Page: page, PerPage: 20},
		)
		return tags, err
	}, func(tag *github.RepositoryTag) {
		tagInfos = append(tagInfos, &platforms.TagInfo{
			TagName: tag.GetName(),
			SHA:     tag.Commit.GetSHA(),
		})
	})
	return tagInfos, err
}

func (p *Platform) GetTagReleaseInfo(ctx context.Context, fullName, tagName string) (*platforms.ReleaseInfo, error) {
	owner, repo, err := x.RepoParseFullName(fullName)
	if err != nil {
		return nil, err
	}
	client := p.GetClient(ctx)
	tag, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tagName)
	if err != nil {
		return nil, err
	}
	tI := &platforms.ReleaseInfo{
		ID:          tag.GetID(),
		TagName:     tag.GetTagName(),
		Title:       tag.GetName(),
		Description: tag.GetBody(),
		FullName:    fullName,
	}
	for _, asset := range tag.Assets {
		tI.Assets = append(tI.Assets, &platforms.AssetInfo{
			ID:   asset.GetID(),
			Name: asset.GetName(),
			URL:  asset.GetBrowserDownloadURL(),
		})
	}
	return tI, nil

}

func (p *Platform) CreateRelease(ctx context.Context, fullName string, releaseInfo *platforms.ReleaseInfo) (*platforms.ReleaseInfo, error) {
	releaseInfo.Init()
	owner, repo, err := x.RepoParseFullName(fullName)
	if err != nil {
		return nil, err
	}
	client := p.GetClient(ctx)
	cRelease, _, err := client.Repositories.CreateRelease(ctx, owner, repo, &github.RepositoryRelease{
		TagName: &releaseInfo.TagName,
		Name:    &releaseInfo.Title,
		Body:    &releaseInfo.Description,
	})
	if err != nil {
		return nil, err
	}
	return &platforms.ReleaseInfo{
		ID:          cRelease.GetID(),
		TagName:     cRelease.GetTagName(),
		Title:       cRelease.GetName(),
		Description: cRelease.GetBody(),
	}, nil
}

func (p *Platform) DeleteReleaseAssets(ctx context.Context, repoInfo *platforms.ReleaseInfo, filenames []string) error {
	owner, repo, err := repoInfo.GetOwnerRepo()
	if err != nil {
		return err
	}
	targetNames := make(map[string]struct{})
	for _, filename := range filenames {
		name := filepath.Base(filename)
		targetNames[name] = struct{}{}
	}
	client := p.GetClient(ctx)
	for _, asset := range repoInfo.Assets {
		if _, ok := targetNames[asset.Name]; ok {
			if _, err = client.Repositories.DeleteReleaseAsset(ctx, owner, repo, asset.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Platform) UploadReleaseAsset(ctx context.Context, repoInfo *platforms.ReleaseInfo, filenames []string) error {
	owner, repo, err := repoInfo.GetOwnerRepo()
	if err != nil {
		return err
	}
	client := p.GetClient(ctx)
	// 遍历待上传文件，逐个上传
	for _, filePath := range filenames {
		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("open file %s failed: %w", filePath, err)
		}
		// ⚠️ 注意：你这里的 defer 在 for 循环里，会等函数退出才执行，可能导致多个文件同时保持打开状态。
		// 更好的做法是显式关闭，或者在上传完毕后立即关闭。
		// defer f.Close()
		// 从路径中提取文件名作为上传名
		assetName := path.Base(filePath)
		// 上传 Release 资产文件
		_, _, err = client.Repositories.UploadReleaseAsset(ctx, owner, repo, repoInfo.ID, &github.UploadOptions{Name: assetName}, f)
		_ = f.Close()
		if err != nil {
			return fmt.Errorf("upload asset %s failed: %w", assetName, err)
		}
	}
	return nil
}
