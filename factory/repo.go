package factory

import (
	"context"
	"fmt"
	"github.com/chihqiang/mpgrm/flags"
	"github.com/chihqiang/mpgrm/pkg/credential"
	"github.com/chihqiang/mpgrm/pkg/logx"
	"github.com/chihqiang/mpgrm/pkg/platforms"
	"github.com/chihqiang/mpgrm/pkg/x"
	"github.com/urfave/cli/v3"
	"net/url"
	"sync"
	"time"
)

// Repo represents a repository with context, CLI command, URL, platform, credentials, and full name.
type Repo struct {
	ctx context.Context // Context for controlling cancellation and deadlines
	cmd *cli.Command    // CLI command instance

	repoURL    *url.URL               // URL of the repository / 仓库 URL
	platform   platforms.IPlatform    // Platform interface for operations (GitHub, Gitee, Gitea, etc.)
	credential *credential.Credential // Authentication credential for the repository
}

// NewRepo creates a new Repo instance based on the CLI command flags and credentials.
// It initializes the repository URL, authentication credentials, platform, and full repository name.
func NewRepo(ctx context.Context, cmd *cli.Command) (*Repo, error) {
	// Initialize Repo instance with context and command
	rt := &Repo{ctx: ctx, cmd: cmd}

	// Get repository URL and credentials from CLI flags
	repoURL, cred, err := flags.GetFormCredential(cmd, true)
	if err != nil {
		return rt, err
	}
	rt.repoURL = repoURL // Assign repository URL
	rt.credential = cred // Assign authentication credential

	// Determine the platform (e.g., GitHub, Gitee, Gitea) using the repository URL and credential
	platform, err := platforms.GetPlatform(repoURL, cred)
	if err != nil {
		return rt, err
	}
	rt.platform = platform // Assign platform
	// Return the initialized Repo instance
	return rt, nil
}

func (r *Repo) ListRepo() (repo []*platforms.RepoInfo, err error) {
	logx.Info("Starting to list repositories...")
	orgName, err := x.RepoURLParseOrgName(r.repoURL)
	if err != nil {
		return nil, err
	}
	logx.Info("Parsed organization/subpath: %s", orgName)
	logx.Info("Fetching repositories from platform...")
	if orgName != "" {
		logx.Info("Organization/subpath detected: %s, fetching org repositories...", orgName)
		repo, err = r.platform.ListOrgRepo(r.ctx, orgName)
	} else {
		logx.Info("No organization/subpath detected, fetching user repositories...")
		repo, err = r.platform.ListUserRepo(r.ctx)
	}
	if err != nil {
		return nil, err
	}
	logx.Info("Successfully fetched %d repositories", len(repo))
	// 遍历仓库并打印名称
	for _, rInfo := range repo {
		logx.Info("  - %s", rInfo.CloneURL)
	}
	return repo, nil
}

func (r *Repo) CloneRepo() error {
	repo, err := r.ListRepo()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, repo := range repo {
		repo := repo
		wg.Add(1)
		go func(repo *platforms.RepoInfo) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			start := time.Now()
			cloneURL := repo.CloneURL
			r.credential.CloneURL = cloneURL
			git, err := NewCredentialGit(r.cmd, r.credential)
			if err != nil {
				logx.Error(" create Git instance for %s: %v", repo.CloneURL, err)
				return
			}
			if err := git.Clone(); err != nil {
				logx.Error("clone %s: %v", repo.CloneURL, err)
				return
			}
			logx.Info("Repository %s cloned successfully (took %s)", repo.CloneURL, time.Since(start))
		}(repo)
	}
	wg.Wait()
	return nil
}
func (r *Repo) RepoSync() error {
	targetURL, targetCredential, err := flags.GetTargetCredential(r.cmd)
	if err != nil {
		return err
	}
	logx.Info("Target URL parsed: %s", targetURL.String())
	targetPlatform, err := platforms.GetPlatform(targetURL, targetCredential)
	if err != nil {
		return err
	}
	repos, err := r.ListRepo()
	if err != nil {
		return err
	}
	logx.Info("Starting repository sync...")
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, repo := range repos {
		repo := repo
		wg.Add(1)
		go func(repo *platforms.RepoInfo) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			start := time.Now()
			cloneURL := repo.CloneURL
			if err := targetCredential.SetCloneByRepoName(repo.Name); err != nil {
				logx.Error(err.Error())
				return
			}
			r.credential.CloneURL = cloneURL
			logx.Info("Source URL: %s", cloneURL)
			logx.Info("Source Credential %s", r.credential)
			logx.Info("Target URL: %s", targetCredential.CloneURL)
			logx.Info("Target Credential %s", targetCredential)
			doubleCredentialGit, err := NewDoubleCredentialGit(r.cmd, r.credential, targetCredential)
			if err != nil {
				logx.Error("Failed to create Git instance for %s: %v", targetCredential.CloneURL, err)
				return
			}
			logx.Info("Git instance created for %s", targetCredential.CloneURL)
			targetFullName, _ := targetCredential.GetFullName()
			detail, err := targetPlatform.GetRepoDetail(r.ctx, targetFullName)
			if err != nil || detail.ID == 0 {
				logx.Warn("Target repository %s does not exist or cannot be fetched, creating...", targetCredential.CloneURL)
				if createErr := targetPlatform.CreateRepo(r.ctx, &platforms.RepoInfo{
					Name:        repo.Name,
					IsPrivate:   repo.IsPrivate,
					FullName:    targetFullName,
					Description: repo.Description,
					Homepage:    repo.Homepage,
				}); createErr != nil {
					logx.Error("create target repository %s: %v", targetCredential.CloneURL, createErr)
					return
				}
			}
			logx.Info("Pushing repository: %s", targetCredential.CloneURL)
			if err := doubleCredentialGit.Push(); err != nil {
				logx.Error("push %s: %v", targetCredential.CloneURL, err)
				return
			}
			logx.Info("Repository %s synced successfully (took %s)", targetCredential.CloneURL, time.Since(start))
		}(repo)
	}
	wg.Wait()
	logx.Info("All repositories sync completed")
	return nil
}

func (r *Repo) CreateRelease() error {
	fullName, err := r.credential.GetFullName()
	if err != nil {
		return err
	}
	tags, err := r.platform.ListTags(r.ctx, fullName)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	for i, tag := range tags {
		start := time.Now()
		_, err := r.platform.CreateRelease(r.ctx, fullName, &platforms.ReleaseInfo{TagName: tag.TagName})
		if err != nil {
			logx.Warn("failed to create release for tag '%s': %v", tag.TagName, err)
			continue
		}
		logx.Info("Created release for tag '%s' (%d/%d) in %s", tag.TagName, i+1, len(tags), time.Since(start))
	}
	return nil
}

func (r *Repo) Upload(tag string, filenames []string) error {
	fullName, err := r.credential.GetFullName()
	if err != nil {
		return err
	}
	logx.Info("Uploading files for tag '%s'...", tag)
	info, err := r.platform.GetTagReleaseInfo(r.ctx, fullName, tag)
	if err != nil {
		return fmt.Errorf("failed to get release info for tag '%s': %w", tag, err)
	}

	if err := r.platform.DeleteReleaseAssets(r.ctx, info, filenames); err != nil {
		logx.Warn("failed to delete existing assets for tag '%s': %v", tag, err)
	}

	if err := r.platform.UploadReleaseAsset(r.ctx, info, filenames); err != nil {
		return fmt.Errorf("failed to upload assets for tag '%s': %w", tag, err)
	}

	logx.Info("Upload completed for tag '%s', %d files", tag, len(filenames))
	return nil
}

func (r *Repo) getReleasePath() (string, error) {
	return r.credential.GetCategoryNamWorkspace(credential.WorkspaceCategoryReleases, flags.GetWorkspace(r.cmd))
}

func (r *Repo) Download(tags []string) (map[string][]string, error) {
	fullName, err := r.credential.GetFullName()
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, fmt.Errorf("no tags provided for download")
	}

	workspace, err := r.getReleasePath()
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}
	logx.Info("Downloading releases into workspace: %s", workspace)
	tagFiles := make(map[string][]string)
	var successCount, failCount int
	var failedTags []string

	for i, tag := range tags {
		start := time.Now()
		logx.Info("Processing tag '%s' (%d/%d)", tag, i+1, len(tags))

		info, err := r.platform.GetTagReleaseInfo(r.ctx, fullName, tag)
		if err != nil {
			logx.Warn("failed to get release info for tag '%s': %v", tag, err)
			failCount++
			failedTags = append(failedTags, tag)
			continue
		}
		files, err := info.Download(workspace)
		if err != nil {
			logx.Warn("failed to download files for tag '%s': %v", tag, err)
			failCount++
			failedTags = append(failedTags, tag)
			continue
		}
		tagFiles[tag] = files
		successCount++
		logx.Info("Downloaded %d files for tag '%s' in %s", len(files), tag, time.Since(start))
		for _, f := range files {
			logx.Info("  - %s", f)
		}
	}
	if failCount > 0 {
		logx.Warn("Download completed: %d success, %d failed (%v)", successCount, failCount, failedTags)
	} else {
		logx.Info("Download completed: %d success, %d failed", successCount, failCount)
	}
	return tagFiles, nil
}
