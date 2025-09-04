package gitea

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"wangzhiqiang/mpgrm/pkg/httpx"
	"wangzhiqiang/mpgrm/pkg/platforms"
	"wangzhiqiang/mpgrm/pkg/x"
)

func (p *Platform) ListOrgRepo(ctx context.Context, orgName string) ([]*platforms.RepoInfo, error) {
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	var rInfos []*platforms.RepoInfo
	err = httpx.Paginate[*gitea.Repository](func(page int) ([]*gitea.Repository, error) {
		repos, _, err := client.ListOrgRepos(orgName, gitea.ListOrgReposOptions{
			ListOptions: gitea.ListOptions{
				Page:     page,
				PageSize: 20,
			},
		})
		return repos, err
	}, func(repo *gitea.Repository) {
		rInfos = append(rInfos, &platforms.RepoInfo{
			ID:          repo.ID,
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: repo.Description,
			Homepage:    repo.Website,
			IsPrivate:   repo.Private,
			CloneURL:    repo.CloneURL,
		})
	})
	return rInfos, err
}

func (p *Platform) ListUserRepo(ctx context.Context) ([]*platforms.RepoInfo, error) {
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	var rInfos []*platforms.RepoInfo
	err = httpx.Paginate[*gitea.Repository](func(page int) ([]*gitea.Repository, error) {
		repos, _, err := client.ListUserRepos(p.Credential.Username, gitea.ListReposOptions{
			ListOptions: gitea.ListOptions{
				Page:     page,
				PageSize: 20,
			},
		})
		return repos, err
	}, func(repo *gitea.Repository) {
		rInfos = append(rInfos, &platforms.RepoInfo{
			ID:          repo.ID,
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: repo.Description,
			Homepage:    repo.Website,
			IsPrivate:   repo.Private,
			CloneURL:    repo.CloneURL,
		})
	})
	return rInfos, err
}

func (p *Platform) GetRepoDetail(ctx context.Context, fullName string) (*platforms.RepoInfo, error) {
	owner, repo, err := x.RepoParseFullName(fullName)
	if err != nil {
		return nil, err
	}
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	repoResp, _, err := client.GetRepo(owner, repo)
	if err != nil {
		return nil, err
	}
	return &platforms.RepoInfo{
		ID:          repoResp.ID,
		Name:        repoResp.Name,
		FullName:    repoResp.FullName,
		Description: repoResp.Description,
		Homepage:    repoResp.Website,
		IsPrivate:   repoResp.Private,
		CloneURL:    repoResp.CloneURL,
	}, nil
}

func (p *Platform) CreateRepo(ctx context.Context, repoInfo *platforms.RepoInfo) error {
	client, err := p.GetClient()
	if err != nil {
		return err
	}
	opt := gitea.CreateRepoOption{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		Private:     repoInfo.IsPrivate,
	}
	orgName, _ := repoInfo.GetOrgName()
	if orgName == "" || orgName == p.Credential.Username {
		_, _, err = client.CreateRepo(opt)
	} else {
		_, _, err = client.CreateOrgRepo(orgName, opt)
	}
	return err
}

func (p *Platform) DeleteRepo(ctx context.Context, repoInfo *platforms.RepoInfo) error {
	owner, repo, err := repoInfo.GetOwnerRepo()
	if err != nil {
		return err
	}
	client, err := p.GetClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteRepo(owner, repo)
	return err
}
