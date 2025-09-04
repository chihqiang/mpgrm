package gitee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"wangzhiqiang/mpgrm/pkg/httpx"
	"wangzhiqiang/mpgrm/pkg/platforms"
)

func (p *Platform) ListOrgRepo(ctx context.Context, orgName string) ([]*platforms.RepoInfo, error) {
	var rInfos []*platforms.RepoInfo
	err := httpx.Paginate[*RepoInfoResponse](func(page int) ([]*RepoInfoResponse, error) {
		var body []*RepoInfoResponse
		_, err := httpx.GetD(ctx, p.GetURLWithToken(fmt.Sprintf("orgs/%s/repos", orgName), map[string]string{
			"page":     fmt.Sprintf("%d", page),
			"per_page": "10",
		}), &body)
		if err != nil {
			return nil, err
		}
		return body, err
	}, func(repo *RepoInfoResponse) {
		rInfos = append(rInfos, &platforms.RepoInfo{
			ID:          repo.Id,
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: repo.Description,
			Homepage:    repo.Homepage,
			IsPrivate:   repo.Private,
			CloneURL:    repo.HtmlUrl,
		})
	})
	return rInfos, err
}

func (p *Platform) ListUserRepo(ctx context.Context) ([]*platforms.RepoInfo, error) {
	var rInfos []*platforms.RepoInfo
	err := httpx.Paginate[*RepoInfoResponse](func(page int) ([]*RepoInfoResponse, error) {
		var body []*RepoInfoResponse
		urlApi := p.GetURLWithToken("user/repos", map[string]string{
			"page":     fmt.Sprintf("%d", page),
			"per_page": "10",
			//"type":     "personal",
		})
		_, err := httpx.GetD(ctx, urlApi, &body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}, func(repo *RepoInfoResponse) {
		rInfos = append(rInfos, &platforms.RepoInfo{
			ID:          repo.Id,
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: repo.Description,
			Homepage:    repo.Homepage,
			IsPrivate:   repo.Private,
			CloneURL:    repo.HtmlUrl,
		})
	})
	return rInfos, err
}

func (p *Platform) GetRepoDetail(ctx context.Context, fullName string) (*platforms.RepoInfo, error) {
	apiURL := p.GetURLWithToken(fmt.Sprintf("repos/%s", fullName), map[string]string{})
	var body RepoInfoResponse
	_, err := httpx.GetD(ctx, apiURL, &body)
	if err != nil {
		return nil, err
	}
	return &platforms.RepoInfo{
		ID:          body.Id,
		Name:        body.Name,
		FullName:    body.FullName,
		Description: body.Description,
		Homepage:    body.Homepage,
		IsPrivate:   body.Private,
	}, nil
}

func (p *Platform) CreateRepo(ctx context.Context, repoInfo *platforms.RepoInfo) error {
	req := struct {
		AccessToken string `json:"access_token"`
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		Homepage    string `json:"homepage,omitempty"`
		Private     bool   `json:"private,omitempty"`
		HasIssues   bool   `json:"has_issues,omitempty"`
		HasWiki     bool   `json:"has_wiki,omitempty"`
		CanComment  bool   `json:"can_comment,omitempty"`
	}{
		AccessToken: p.Credential.Token,
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		Homepage:    repoInfo.Homepage,
		Private:     repoInfo.IsPrivate,
		HasIssues:   true,
		HasWiki:     true,
		CanComment:  true,
	}
	orgName, _ := repoInfo.GetOrgName()
	var apiUrl string
	if orgName == "" || orgName == p.Credential.Username {
		apiUrl = p.GetURLWithToken("user/repos", map[string]string{})
	} else {
		apiUrl = p.GetURLWithToken(fmt.Sprintf("orgs/%s/repos", orgName), map[string]string{})
	}
	jsonReq, _ := json.Marshal(req)
	var resp CreateRepoResponse
	_, err := httpx.PostD(ctx, apiUrl, bytes.NewReader(jsonReq), &resp, map[string]string{
		"Content-Type": "application/json",
	})
	return err
}

func (p *Platform) DeleteRepo(ctx context.Context, repoInfo *platforms.RepoInfo) error {
	apiUrl := p.GetURLWithToken(fmt.Sprintf("repos/%s", repoInfo.FullName), map[string]string{})
	_, err := httpx.Request(context.Background(), http.MethodDelete, apiUrl, nil, map[string]string{})
	return err
}
