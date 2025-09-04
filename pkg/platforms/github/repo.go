package github

import (
	"context"
	"github.com/google/go-github/v73/github"
	"wangzhiqiang/mpgrm/pkg/httpx"
	"wangzhiqiang/mpgrm/pkg/platforms"
	"wangzhiqiang/mpgrm/pkg/x"
)

func (p *Platform) ListOrgRepo(ctx context.Context, orgName string) ([]*platforms.RepoInfo, error) {
	var rInfos []*platforms.RepoInfo
	client := p.GetClient(ctx)
	err := httpx.Paginate[*github.Repository](func(page int) ([]*github.Repository, error) {
		repos, _, err := client.Repositories.ListByOrg(context.Background(), orgName, &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{
				Page: page,
			},
		})
		return repos, err
	}, func(repo *github.Repository) {
		rInfos = append(rInfos, &platforms.RepoInfo{
			ID:          repo.GetID(),
			Name:        repo.GetName(),
			FullName:    repo.GetFullName(),
			Description: repo.GetDescription(),
			Homepage:    repo.GetHomepage(),
			IsPrivate:   repo.GetPrivate(),
			CloneURL:    repo.GetCloneURL(),
		})
	})
	return rInfos, err
}

func (p *Platform) ListUserRepo(ctx context.Context) ([]*platforms.RepoInfo, error) {
	var rInfos []*platforms.RepoInfo
	client := p.GetClient(ctx)
	err := httpx.Paginate[*github.Repository](func(page int) ([]*github.Repository, error) {
		repos, _, err := client.Repositories.ListByUser(ctx, p.Credential.Username, &github.RepositoryListByUserOptions{
			ListOptions: github.ListOptions{Page: page},
		})
		return repos, err
	}, func(repo *github.Repository) {
		rInfos = append(rInfos, &platforms.RepoInfo{
			ID:          repo.GetID(),
			Name:        repo.GetName(),
			FullName:    repo.GetFullName(),
			Description: repo.GetDescription(),
			Homepage:    repo.GetHomepage(),
			IsPrivate:   repo.GetPrivate(),
			CloneURL:    repo.GetCloneURL(),
		})
	})
	return rInfos, err
}

func (p *Platform) GetRepoDetail(ctx context.Context, fullName string) (*platforms.RepoInfo, error) {
	client := p.GetClient(ctx)
	owner, repo, err := x.RepoParseFullName(fullName)
	if err != nil {
		return nil, err
	}
	repoG, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return nil, err
	}
	ri := &platforms.RepoInfo{
		ID:          repoG.GetID(),
		Name:        repoG.GetName(),
		FullName:    repoG.GetFullName(),
		Description: repoG.GetDescription(),
		Homepage:    repoG.GetHomepage(),
		IsPrivate:   repoG.GetPrivate(),
		CloneURL:    repoG.GetCloneURL(),
	}
	return ri, nil
}

func (p *Platform) CreateRepo(ctx context.Context, repoInfo *platforms.RepoInfo) error {
	owner, _, err := repoInfo.GetOwnerRepo()
	if err != nil {
		return err
	}
	client := p.GetClient(ctx)
	_, _, err = client.Repositories.Create(context.Background(), owner, &github.Repository{
		Name:        &repoInfo.Name,
		Private:     &repoInfo.IsPrivate,
		Description: &repoInfo.Description,
		Homepage:    &repoInfo.Homepage,
	})
	return err
}

func (p *Platform) DeleteRepo(ctx context.Context, repoInfo *platforms.RepoInfo) error {
	client := p.GetClient(ctx)
	owner, repo, err := repoInfo.GetOwnerRepo()
	if err != nil {
		return err
	}
	_, err = client.Repositories.Delete(context.Background(), owner, repo)
	return err
}
