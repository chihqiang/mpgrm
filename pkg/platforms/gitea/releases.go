package gitea

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"wangzhiqiang/mpgrm/pkg/httpx"
	"wangzhiqiang/mpgrm/pkg/platforms"
	"wangzhiqiang/mpgrm/pkg/x"
)

func (p *Platform) ListTags(ctx context.Context, fullName string) ([]*platforms.TagInfo, error) {
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	owner, repo, err := x.RepoParseFullName(fullName)
	if err != nil {
		return nil, err
	}
	var allTags []*platforms.TagInfo
	err = httpx.Paginate[*gitea.Tag](func(page int) ([]*gitea.Tag, error) {
		tags, _, err := client.ListRepoTags(owner, repo, gitea.ListRepoTagsOptions{
			ListOptions: gitea.ListOptions{
				Page:     page,
				PageSize: 20,
			},
		})
		return tags, err
	}, func(tag *gitea.Tag) {
		allTags = append(allTags, &platforms.TagInfo{
			TagName: tag.Name,
			SHA:     tag.Commit.SHA,
		})
	})
	return allTags, err
}

func (p *Platform) GetTagReleaseInfo(ctx context.Context, fullName, tagName string) (*platforms.ReleaseInfo, error) {
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	owner, repo, err := x.RepoParseFullName(fullName)
	if err != nil {
		return nil, err
	}
	release, _, err := client.GetReleaseByTag(owner, repo, tagName)
	if err != nil {
		return nil, fmt.Errorf("failed to get release by tag %s: %w", tagName, err)
	}
	inf := &platforms.ReleaseInfo{
		ID:          release.ID,
		TagName:     release.TagName,
		Title:       release.Title,
		Description: release.Note,

		FullName: fullName,
	}
	for _, attachment := range release.Attachments {
		inf.Assets = append(inf.Assets, &platforms.AssetInfo{
			ID:   attachment.ID,
			Name: attachment.Name,
			URL:  attachment.DownloadURL,
		})
	}
	return inf, nil
}

func (p *Platform) CreateRelease(ctx context.Context, fullName string, releaseInfo *platforms.ReleaseInfo) (newTagInfo *platforms.ReleaseInfo, er error) {
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	releaseInfo.Init()
	owner, repo, err := x.RepoParseFullName(fullName)
	if err != nil {
		return nil, err
	}
	release, _, err := client.CreateRelease(owner, repo, gitea.CreateReleaseOption{
		TagName: releaseInfo.TagName,
		Title:   releaseInfo.Title,
		Note:    releaseInfo.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create release: %w", err)
	}
	return &platforms.ReleaseInfo{
		ID:          release.ID,
		TagName:     release.TagName,
		Title:       release.Title,
		Description: release.Note,
	}, nil
}

func (p *Platform) DeleteReleaseAssets(ctx context.Context, repoInfo *platforms.ReleaseInfo, filenames []string) error {
	client, err := p.GetClient()
	if err != nil {
		return err
	}
	owner, repo, err := repoInfo.GetOwnerRepo()
	if err != nil {
		return err
	}
	targetNames := make(map[string]struct{})
	for _, filename := range filenames {
		name := filepath.Base(filename)
		targetNames[name] = struct{}{}
	}
	for _, asset := range repoInfo.Assets {
		if _, ok := targetNames[asset.Name]; ok {
			if _, err := client.DeleteReleaseAttachment(owner, repo, repoInfo.ID, asset.ID); err != nil {
				return err
			}
		}
	}
	return nil

}

func (p *Platform) UploadReleaseAsset(ctx context.Context, tagInfo *platforms.ReleaseInfo, filenames []string) error {
	client, err := p.GetClient()
	if err != nil {
		return err
	}
	owner, repo, err := tagInfo.GetOwnerRepo()
	if err != nil {
		return err
	}
	// 遍历待上传文件，逐个上传
	for _, filePath := range filenames {
		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("打开文件 %s 失败: %w", filePath, err)
		}
		_, _, err = client.CreateReleaseAttachment(owner, repo, tagInfo.ID, f, path.Base(filePath))
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()
	}
	return nil
}
