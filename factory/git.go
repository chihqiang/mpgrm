package factory

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"net/url"
	"time"
	"wangzhiqiang/mpgrm/flags"
	"wangzhiqiang/mpgrm/pkg/credential"
	"wangzhiqiang/mpgrm/pkg/gitx"
	"wangzhiqiang/mpgrm/pkg/x"
)

type Git struct {
	ctx context.Context

	workspace      string
	branches, tags []string

	credential       *credential.Credential
	fullName         string
	targetCredential *credential.Credential
	targetFullName   string
}

func NewGit(cmd *cli.Command, credential *credential.Credential, targetCredential *credential.Credential) (*Git, error) {
	rt := &Git{
		ctx:              context.Background(),
		workspace:        flags.GetWorkspace(cmd),
		branches:         flags.GetBranches(cmd),
		tags:             flags.GetTags(cmd),
		credential:       credential,
		targetCredential: targetCredential,
	}
	repoURL, err := url.Parse(credential.CloneURL)
	if err != nil {
		return rt, err
	}
	fullName, err := x.RepoURLParseFullName(repoURL)
	if err != nil {
		return rt, err
	}
	rt.fullName = fullName
	targetURL, err := url.Parse(targetCredential.CloneURL)
	if err != nil {
		return rt, err
	}
	targetFullName, err := x.RepoURLParseFullName(targetURL)
	if err != nil {
		return rt, err
	}
	rt.targetFullName = targetFullName
	return rt, nil
}

func NewGitCmd(ctx context.Context, cmd *cli.Command) (*Git, error) {
	rt := &Git{
		ctx:       ctx,
		workspace: flags.GetWorkspace(cmd),
		branches:  flags.GetBranches(cmd),
		tags:      flags.GetTags(cmd),
	}
	// 获取源仓库 URL 和认证
	repoURL, cred, _ := flags.GetFormCredential(cmd, false)
	rt.credential = cred

	// 解析源仓库全名
	fullName, err := x.RepoURLParseFullName(repoURL)
	if err != nil {
		return rt, err
	}
	rt.fullName = fullName

	// 获取目标仓库 URL 和认证
	targetRepoURL, targetCred, err := flags.GetTargetCredential(cmd)
	if err != nil {
		return rt, err
	}
	rt.targetCredential = targetCred
	// 解析目标仓库全名
	targetFullName, err := x.RepoURLParseFullName(targetRepoURL)
	if err != nil {
		return rt, err
	}
	rt.targetFullName = targetFullName
	return rt, nil
}

func (g *Git) Push() error {
	log.Printf("Starting git migration from %s to %s", g.fullName, g.targetFullName)

	start := time.Now()
	migrate := gitx.NewGitMigrate(g.credential, g.targetCredential)

	// 获取 workspace
	workspace, err := g.credential.GetCategoryNamWorkspace(credential.WorkspaceCategoryGit, g.workspace)
	if err != nil {
		log.Printf("Failed to get workspace: %v", err)
		return err
	}
	log.Printf("Using workspace: %s", workspace)

	// 获取分支和标签
	branches := g.branches
	tags := g.tags
	log.Printf("Preparing to clone branches: %v, tags: %v", branches, tags)

	// Clone 仓库并获取实际分支和标签
	actualBranches, actualTags, err := migrate.Clone(workspace, branches, tags)
	if err != nil {
		log.Printf("Failed to clone repository: %v", err)
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	log.Printf("Repository cloned successfully. Cloned branches: %v, tags: %v", actualBranches, actualTags)

	// Push 到目标仓库
	if err := migrate.Push(workspace); err != nil {
		log.Printf("Failed to push repository: %v", err)
		return err
	}

	elapsed := time.Since(start)
	log.Printf("Push to target repository completed successfully, total elapsed time: %s", elapsed)
	return nil
}
