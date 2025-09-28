package cnb

import (
	"bytes"
	"cnb.cool/cnb/sdk/go-cnb/cnb"
	"cnb.cool/cnb/sdk/go-cnb/cnb/types/api"
	"context"
	"fmt"
	"github.com/chihqiang/mpgrm/pkg/httpx"
	"github.com/chihqiang/mpgrm/pkg/platforms"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

func (p *Platform) ListTags(ctx context.Context, fullName string) ([]*platforms.TagInfo, error) {
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	var allTags []*platforms.TagInfo
	err = httpx.Paginate[*api.Tag](func(page int) ([]*api.Tag, error) {
		tags, _, err := client.Git.ListTags(ctx, fullName, &cnb.ListTagsOptions{
			Page:     page,
			PageSize: 10,
		})
		return tags, err
	}, func(tag *api.Tag) {
		allTags = append(allTags, &platforms.TagInfo{
			TagName: tag.Name,
			SHA:     tag.Commit.Sha,
		})
	})
	return nil, err
}

func (p *Platform) GetTagReleaseInfo(ctx context.Context, fullName, tagName string) (*platforms.ReleaseInfo, error) {
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	byTag, _, err := client.Releases.GetReleaseByTag(ctx, fullName, tagName)
	if err != nil {
		return nil, fmt.Errorf("get release by tag failed: %w", err)
	}
	riID, _ := strconv.ParseInt(byTag.Id, 10, 64)
	rfo := &platforms.ReleaseInfo{
		ID:          riID,
		TagName:     byTag.TagName,
		Title:       byTag.Name,
		Description: byTag.Body,
		FullName:    fullName,
	}
	for _, asset := range byTag.Assets {
		aID, _ := strconv.ParseInt(asset.Id, 10, 64)
		rfo.Assets = append(rfo.Assets, &platforms.AssetInfo{
			ID:   aID,
			Name: asset.Name,
			URL:  fmt.Sprintf("%s%s", downloadURL, asset.Path),
		})
	}
	return rfo, nil
}

func (p *Platform) CreateRelease(ctx context.Context, fullName string, releaseInfo *platforms.ReleaseInfo) (newTagInfo *platforms.ReleaseInfo, er error) {
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	releaseInfo.Init()
	release, _, err := client.Releases.PostRelease(ctx, fullName, &cnb.PostReleaseRequest{
		Name:    releaseInfo.Title,
		TagName: releaseInfo.TagName,
		Body:    releaseInfo.Description,
	})
	if err != nil {
		return nil, err
	}
	rID, _ := strconv.ParseInt(release.Id, 10, 64)
	return &platforms.ReleaseInfo{
		ID:          rID,
		TagName:     release.TagName,
		Title:       release.Name,
		Description: release.Body,
	}, nil
}

func (p *Platform) DeleteReleaseAssets(ctx context.Context, repoInfo *platforms.ReleaseInfo, filenames []string) error {
	client, err := p.GetClient()
	if err != nil {
		return err
	}
	targetNames := make(map[string]struct{})
	for _, path := range filenames {
		name := filepath.Base(path)
		targetNames[name] = struct{}{}
	}
	for _, asset := range repoInfo.Assets {
		if _, ok := targetNames[asset.Name]; ok {
			releaseID := strconv.FormatInt(repoInfo.ID, 10)
			assetID := strconv.FormatInt(asset.ID, 10)
			if _, err := client.Releases.DeleteReleaseAsset(context.Background(), repoInfo.FullName, releaseID, assetID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Platform) UploadReleaseAsset(ctx context.Context, tagInfo *platforms.ReleaseInfo, filenames []string) error {
	for _, file := range filenames {
		tagID := strconv.FormatInt(tagInfo.ID, 10)
		if err := p.uploadAsset(tagInfo.FullName, file, tagID); err != nil {
			return err
		}
	}
	return nil
}

func (p *Platform) uploadAsset(repo, filename string, releaseID string) error {
	client, err := p.GetClient()
	if err != nil {
		return err
	}
	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("os.Stat error: %v", err)
	}
	AssetName := filepath.Base(filename)
	r, _, err := client.Releases.PostReleaseAssetUploadURL(context.Background(), repo, releaseID, &cnb.PostReleaseAssetUploadURLRequest{
		AssetName: AssetName,
		Size:      int(info.Size()),
		Overwrite: true,
	})
	if err != nil {
		return err
	}
	if err = PutObjectToCos(r.UploadUrl, filename); err != nil {
		return err
	}
	return AssetUploadConfirmation(p.Credential.Token, r.VerifyUrl)
}

func PutObjectToCos(url string, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("put object response code: %d, detail: %s", resp.StatusCode, string(body))
	}
	return nil
}
func AssetUploadConfirmation(token, verifyUrl string) error {
	decodedURL, err := url.QueryUnescape(verifyUrl)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, decodedURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.cnb.api+json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("verify asset response code: %d, detail: %s", resp.StatusCode, string(body))
	}
	return nil
}
