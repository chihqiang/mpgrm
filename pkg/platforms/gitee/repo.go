package gitee

import (
	"cnb.cool/zhiqiangwang/pkg/go-gitee/gitee"
	"cnb.cool/zhiqiangwang/pkg/go-gitee/gitee/types/ibase"
	"context"
	"wangzhiqiang/mpgrm/pkg/httpx"
	"wangzhiqiang/mpgrm/pkg/platforms"
	"wangzhiqiang/mpgrm/pkg/x"
)

func (p *Platform) ListOrgRepo(ctx context.Context, orgName string) ([]*platforms.RepoInfo, error) {
	var rInfos []*platforms.RepoInfo
	err := httpx.Paginate[*ibase.Project](func(page int) ([]*ibase.Project, error) {
		client := p.Client()
		repos, _, err := client.Repositories.GetV5OrgsOrgRepos(ctx, orgName, &gitee.GetV5OrgsOrgReposOptions{
			Page: page,
		})
		return repos, err
	}, func(repo *ibase.Project) {
		rInfos = append(rInfos, &platforms.RepoInfo{
			ID:          int64(repo.Id),
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
	err := httpx.Paginate[*ibase.Project](func(page int) ([]*ibase.Project, error) {
		client := p.Client()
		repos, _, err := client.Repositories.GetV5UserRepos(ctx, &gitee.GetV5UserReposOptions{
			Page: page,
		})
		return repos, err
	}, func(repo *ibase.Project) {
		rInfos = append(rInfos, &platforms.RepoInfo{
			ID:          int64(repo.Id),
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
	client := p.Client()
	owner, repo, err := x.RepoParseFullName(fullName)
	if err != nil {
		return nil, err
	}
	body, _, err := client.Repositories.GetV5ReposOwnerRepo(ctx, owner, repo, &gitee.GetV5ReposOwnerRepoOptions{})
	if err != nil {
		return nil, err
	}
	return &platforms.RepoInfo{
		ID:          int64(body.Id),
		Name:        body.Name,
		FullName:    body.FullName,
		Description: body.Description,
		Homepage:    body.Homepage,
		IsPrivate:   body.Private,
		CloneURL:    body.HtmlUrl,
	}, nil
}

func (p *Platform) CreateRepo(ctx context.Context, repoInfo *platforms.RepoInfo) error {
	orgName, _ := repoInfo.GetOrgName()
	client := p.Client()
	var err error
	if orgName == "" || orgName == p.Credential.Username {
		_, _, err = client.Repositories.PostV5UserRepos(ctx, &gitee.PostV5UserReposForms{
			Name:        repoInfo.Name,
			Description: repoInfo.Description,
			Homepage:    repoInfo.Homepage,
			Private:     repoInfo.IsPrivate,
			HasIssues:   true,
			HasWiki:     true,
			CanComment:  true,
		})
	} else {
		_, _, err = client.Repositories.PostV5OrgsOrgRepos(ctx, orgName, &gitee.PostV5OrgsOrgReposForms{
			Name:        repoInfo.Name,
			Description: repoInfo.Description,
			Homepage:    repoInfo.Homepage,
			Private:     repoInfo.IsPrivate,
			HasIssues:   true,
			HasWiki:     true,
			CanComment:  true,
		})
	}
	return err
}

func (p *Platform) DeleteRepo(ctx context.Context, repoInfo *platforms.RepoInfo) error {
	client := p.Client()
	owner, repo, err := repoInfo.GetOwnerRepo()
	if err != nil {
		return err
	}
	_, err = client.Repositories.DeleteV5ReposOwnerRepo(ctx, owner, repo, &gitee.DeleteV5ReposOwnerRepoOptions{})
	return err
}
