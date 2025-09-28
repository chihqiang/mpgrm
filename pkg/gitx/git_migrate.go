package gitx

import (
	"errors"
	"fmt"
	"github.com/chihqiang/mpgrm/pkg/credential"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitMigrate struct {
	form   *credential.Credential
	target *credential.Credential
}

// NewGitMigrateDouble 创建 GitMigrate 实例并初始化认证
func NewGitMigrateDouble(from, target *credential.Credential) *GitMigrate {
	m := NewGitMigrate()
	m.WithForm(from)
	m.WithTarget(target)
	return m
}
func NewGitMigrate() *GitMigrate {
	m := &GitMigrate{}
	return m
}

func (m *GitMigrate) WithForm(credential *credential.Credential) {
	m.form = credential
}

func (m *GitMigrate) WithTarget(credential *credential.Credential) {
	m.target = credential
}
func (m *GitMigrate) Push(path string) error {
	// 打开本地仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed target open repo: %w", err)
	}
	// 设置远程 "target"
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "target",
		URLs: []string{m.target.CloneURL},
	})
	if err != nil && err != git.ErrRemoteExists {
		return fmt.Errorf("failed target create remote: %w", err)
	}
	// 推送所有分支
	err = repo.Push(&git.PushOptions{
		RemoteName: "target",
		RefSpecs: []config.RefSpec{
			"+refs/heads/*:refs/heads/*", // 推送所有分支
		},
		Auth:  m.target.GetGitAuth(),
		Force: true,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed target push branches: %w", err)
	}
	// 推送所有 tag
	err = repo.Push(&git.PushOptions{
		RemoteName: "target",
		RefSpecs: []config.RefSpec{
			"+refs/tags/*:refs/tags/*", // 推送所有 tag
		},
		Auth: m.target.GetGitAuth(),
	})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return nil
		}
		return fmt.Errorf("failed target push tags: %w", err)
	}
	return nil
}

func (m *GitMigrate) Clone(path string, branches, tags []string) ([]string, []string, error) {
	var (
		err            error
		rBranch, rTags []string
		refSpecs       []config.RefSpec
	)

	// 如果没有指定分支或标签，从远程获取
	if len(branches) == 0 || len(tags) == 0 {
		rBranch, rTags, err = m.getRemoteBranchAndTag()
		if err != nil {
			return nil, nil, fmt.Errorf("failed target get remote branches and tags: %w", err)
		}
	}
	if len(branches) == 0 {
		branches = rBranch
	}
	if len(tags) == 0 {
		tags = rTags
	}

	// 构建 RefSpecs
	for _, branch := range branches {
		if branch == "" {
			continue
		}
		refSpecs = append(refSpecs, config.RefSpec(
			fmt.Sprintf("+refs/heads/%s:refs/remotes/origin/%s", branch, branch),
		))
	}
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		refSpecs = append(refSpecs, config.RefSpec(
			fmt.Sprintf("+refs/tags/%s:refs/tags/%s", tag, tag),
		))
	}

	if len(refSpecs) == 0 {
		return nil, nil, fmt.Errorf("no remote branches or tags specified")
	}

	// 初始化空仓库
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return nil, nil, fmt.Errorf("failed target init repo: %w", err)
	}

	// 添加远程
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{m.form.CloneURL},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed target create remote: %w", err)
	}

	// Fetch 指定 RefSpecs
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   refSpecs,
		Auth:       m.form.GetGitAuth(),
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, nil, fmt.Errorf("failed target fetch refs: %w", err)
	}

	// 获取工作区
	w, err := repo.Worktree()
	if err != nil {
		return nil, nil, fmt.Errorf("failed target get worktree: %w", err)
	}

	// 遍历所有分支，创建本地分支并 checkout
	for _, branch := range branches {
		if branch == "" {
			continue
		}
		remoteRef := plumbing.NewRemoteReferenceName("origin", branch)
		ref, err := repo.Reference(remoteRef, true)
		if err != nil {
			return nil, nil, fmt.Errorf("remote branch not found: %s", branch)
		}

		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branch),
			Create: true,
			Force:  true,
			Hash:   ref.Hash(),
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed target checkout branch %s: %w", branch, err)
		}
	}

	// 遍历标签，创建本地 tag
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		remoteTag := plumbing.NewTagReferenceName(tag)
		ref, err := repo.Reference(remoteTag, true)
		if err != nil {
			return nil, nil, fmt.Errorf("remote tag not found: %s", tag)
		}
		// 创建本地 tag
		_, err = repo.CreateTag(tag, ref.Hash(), nil)
		if err != nil && !errors.Is(err, git.ErrTagExists) {
			return nil, nil, fmt.Errorf("failed target create tag %s: %w", tag, err)
		}
	}

	// 返回实际使用的分支和标签
	return branches, tags, nil
}

// getRemoteBranches 获取远程分支列表
func (m *GitMigrate) getRemoteBranchAndTag() (branches []string, tags []string, err error) {
	remote := git.NewRemote(nil, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{m.form.CloneURL},
	})
	listOpts := &git.ListOptions{
		Auth: m.form.GetGitAuth(),
	}
	refs, err := remote.List(listOpts)
	if err != nil {
		return []string{}, []string{}, err
	}
	for _, ref := range refs {
		if ref.Name().IsBranch() {
			branches = append(branches, ref.Name().Short())
		}
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}
	return branches, tags, nil
}
