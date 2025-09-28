package cnb

import (
	"cnb.cool/cnb/sdk/go-cnb/cnb"
	"cnb.cool/cnb/sdk/go-cnb/cnb/types/dto"
	"context"
	"fmt"
	"github.com/chihqiang/mpgrm/pkg/httpx"
	"github.com/chihqiang/mpgrm/pkg/platforms"
	"strconv"
)

const VisibilityLevelSecret = "Secret"
const VisibilityLevelPrivate = "Private"

func (p *Platform) ListOrgRepo(ctx context.Context, orgName string) ([]*platforms.RepoInfo, error) {
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	var rInfos []*platforms.RepoInfo
	err = httpx.Paginate[*dto.Repos4User](func(page int) ([]*dto.Repos4User, error) {
		repos, _, err := client.Repositories.GetGroupSubRepos(ctx, orgName, &cnb.GetGroupSubReposOptions{
			Page:     page,
			PageSize: 10,
		})
		return repos, err
	}, func(repo *dto.Repos4User) {
		if repo.VisibilityLevel == VisibilityLevelSecret {
			return
		}
		idInt64, _ := strconv.ParseInt(repo.Id, 10, 64)
		rInfo := &platforms.RepoInfo{
			ID:          idInt64,
			Name:        repo.Name,
			FullName:    repo.Path,
			Description: repo.Description,
			Homepage:    repo.Site,
			IsPrivate:   repo.VisibilityLevel == VisibilityLevelPrivate,
		}
		rInfo.CloneURL = fmt.Sprintf("%s/%s.git", cloneURL, rInfo.FullName)
		rInfos = append(rInfos, rInfo)
	})
	return rInfos, err
}

func (p *Platform) ListUserRepo(ctx context.Context) ([]*platforms.RepoInfo, error) {
	return p.ListOrgRepo(ctx, p.Credential.Username)
}

func (p *Platform) GetRepoDetail(ctx context.Context, fullName string) (*platforms.RepoInfo, error) {
	client, err := p.GetClient()
	if err != nil {
		return nil, err
	}
	repo, _, err := client.Repositories.GetByID(ctx, fullName)
	if err != nil {
		return nil, err
	}
	//https://cnb.cool/cnb/sdk/go-cnb/-/issues/55
	idInt64, _ := strconv.ParseInt(repo.Id, 10, 64)
	return &platforms.RepoInfo{
		ID:          idInt64,
		Name:        repo.Name,
		FullName:    repo.Path,
		Description: repo.Description,
		Homepage:    repo.Site,
		IsPrivate:   repo.VisibilityLevel == VisibilityLevelPrivate,
	}, nil
}

func (p *Platform) CreateRepo(ctx context.Context, repo *platforms.RepoInfo) error {
	client, err := p.GetClient()
	if err != nil {
		return err
	}
	orgName, err := repo.GetOrgName()
	if err != nil {
		return err
	}
	//https://cnb.cool/cnb/sdk/go-cnb/-/issues/57
	var visibility dto.CreateRepoReqVisibility
	if repo.IsPrivate {
		visibility = dto.CreateRepoReqVisibilityPrivate
	} else {
		visibility = dto.CreateRepoReqVisibilityPublic
	}
	_, err = client.Repositories.CreateRepo(ctx, orgName, &cnb.CreateRepoRequest{
		Name:        repo.Name,
		Description: repo.Description,
		Visibility:  visibility,
	})
	return err
}

func (p *Platform) DeleteRepo(ctx context.Context, repo *platforms.RepoInfo) error {
	client, err := p.GetClient()
	if err != nil {
		return err
	}
	_, err = client.Repositories.DeleteRepo(ctx, repo.FullName)
	return err
}
