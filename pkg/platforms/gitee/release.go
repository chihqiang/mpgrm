package gitee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"wangzhiqiang/mpgrm/pkg/httpx"
	"wangzhiqiang/mpgrm/pkg/platforms"
)

func (p *Platform) ListTags(ctx context.Context, fullName string) ([]*platforms.TagInfo, error) {
	var (
		allTags []*platforms.TagInfo
	)
	err := httpx.Paginate[TagResponse](func(page int) ([]TagResponse, error) {
		var result []TagResponse
		apiUrl := p.GetURLWithToken(fmt.Sprintf("/repos/%s/tags", fullName), map[string]string{
			"page":     fmt.Sprintf("%d", page),
			"per_page": "10",
		})
		_, err := httpx.GetD(ctx, apiUrl, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}, func(item TagResponse) {
		allTags = append(allTags, &platforms.TagInfo{
			TagName: item.Name,
			SHA:     item.Commit.Sha,
		})
	})
	return allTags, err
}

func (p *Platform) GetTagReleaseInfo(ctx context.Context, fullName, tagName string) (*platforms.ReleaseInfo, error) {
	apiUrl := p.GetURLWithToken(fmt.Sprintf("repos/%s/releases/tags/%s", fullName, tagName), map[string]string{})
	var releasesTagResponse ReleasesTagResponse
	_, err := httpx.GetD(ctx, apiUrl, &releasesTagResponse)
	if err != nil {
		return nil, fmt.Errorf("request release by tag failed: %w", err)
	}
	if releasesTagResponse.TagName != tagName {
		return nil, fmt.Errorf("release by tag %s not found", tagName)
	}
	info := &platforms.ReleaseInfo{
		ID:          releasesTagResponse.ID,
		TagName:     releasesTagResponse.TagName,
		Title:       releasesTagResponse.Name,
		Description: releasesTagResponse.Body,
		FullName:    fullName,
	}
	attApiUrl := p.GetURLWithToken(fmt.Sprintf("repos/%s/releases/%d/attach_files", fullName, info.ID), map[string]string{})
	var listReleaseAssetResponse []ListReleaseAssetResponse
	_, err = httpx.GetD(ctx, attApiUrl, &listReleaseAssetResponse)
	if err == nil {
		for _, release := range listReleaseAssetResponse {
			if strings.Contains(release.BrowserDownloadUrl, "archive/refs/tags/") {
				continue
			}
			info.Assets = append(info.Assets, &platforms.AssetInfo{
				ID:   release.ID,
				Name: release.Name,
				URL:  release.BrowserDownloadUrl,
			})
		}
	}
	return info, nil
}

func (p *Platform) CreateRelease(ctx context.Context, fullName string, releaseInfo *platforms.ReleaseInfo) (newTagInfo *platforms.ReleaseInfo, er error) {
	releaseInfo.Init()
	jsonData, _ := json.Marshal(map[string]interface{}{
		"tag_name":         releaseInfo.TagName,
		"name":             releaseInfo.Title,
		"body":             releaseInfo.Description,
		"target_commitish": releaseInfo.TagName,
	})
	apiURL := p.GetURLWithToken(fmt.Sprintf("repos/%s/releases", fullName), map[string]string{})
	var result CreateReleaseResponse
	if _, err := httpx.PostD(ctx, apiURL, bytes.NewBuffer(jsonData), &result, map[string]string{
		"content-type": "application/json;charset=UTF-8",
	}); err != nil {
		return nil, err
	}
	inf := &platforms.ReleaseInfo{
		ID:          result.ID,
		TagName:     result.TagName,
		Title:       result.Name,
		Description: result.Body,
	}
	return inf, nil
}

func (p *Platform) DeleteReleaseAssets(ctx context.Context, releaseInfo *platforms.ReleaseInfo, filenames []string) error {
	targetNames := make(map[string]struct{})
	for _, filename := range filenames {
		name := filepath.Base(filename)
		targetNames[name] = struct{}{}
	}
	for _, asset := range releaseInfo.Assets {
		if _, ok := targetNames[asset.Name]; ok {
			apiURL := p.GetURLWithToken(fmt.Sprintf("repos/%s/releases/%d/attach_files/%d", releaseInfo.FullName, releaseInfo.ID, asset.ID), map[string]string{})
			_, err := httpx.Request(ctx, http.MethodDelete, apiURL, nil, map[string]string{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Platform) UploadReleaseAsset(ctx context.Context, releaseInfo *platforms.ReleaseInfo, filenames []string) error {
	for _, filename := range filenames {
		if err := p.uploadAttach(ctx, releaseInfo.FullName, releaseInfo.ID, filename); err != nil {
			return fmt.Errorf("upload attachments %s err: %w", filename, err)
		}
	}
	return nil
}
func (p *Platform) uploadAttach(ctx context.Context, fullName string, releaseID int64, file string) error {
	uploadURL := p.GetURLWithToken(fmt.Sprintf("repos/%s/releases/%d/attach_files", fullName, releaseID), map[string]string{})
	_, err := httpx.Upload(ctx, uploadURL, file, "file", map[string]string{})
	return err
}
