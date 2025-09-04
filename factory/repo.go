package factory

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"net/url"
	"sync"
	"time"
	"wangzhiqiang/mpgrm/flags"
	"wangzhiqiang/mpgrm/pkg/credential"
	"wangzhiqiang/mpgrm/pkg/platforms"
	"wangzhiqiang/mpgrm/pkg/x"
)

// Repo represents a repository with context, CLI command, URL, platform, credentials, and full name.
type Repo struct {
	ctx context.Context // Context for controlling cancellation and deadlines
	cmd *cli.Command    // CLI command instance

	repoURL    *url.URL               // URL of the repository / 仓库 URL
	platform   platforms.IPlatform    // Platform interface for operations (GitHub, Gitee, Gitea, etc.)
	credential *credential.Credential // Authentication credential for the repository
	fullName   string                 // Full repository name in format "owner/repo[/subdir]"
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

	// Parse the full name of the repository (format: owner/repo or owner/repo/subdir)
	fullName, err := cred.GetFullName()
	if err != nil {
		return rt, err
	}
	rt.fullName = fullName // Assign full repository name

	// Return the initialized Repo instance
	return rt, nil
}
func (r *Repo) ListRepo() error {
	log.Println("Starting to list repositories...")
	orgName, err := x.RepoURLParseOrgName(r.repoURL)
	if err != nil {
		return err
	}
	log.Printf("Parsed organization/subpath: %s", orgName)

	var repo []*platforms.RepoInfo

	log.Println("Fetching repositories from platform...")
	if orgName != "" {
		log.Printf("Organization/subpath detected: %s, fetching org repositories...", orgName)
		repo, err = r.platform.ListOrgRepo(r.ctx, orgName)
	} else {
		log.Println("No organization/subpath detected, fetching user repositories...")
		repo, err = r.platform.ListUserRepo(r.ctx)
	}
	if err != nil {
		return err
	}
	log.Printf("Successfully fetched %d repositories", len(repo))
	// 遍历仓库并打印名称
	for _, rInfo := range repo {
		log.Printf("  - %s", rInfo.CloneURL)
	}
	return nil
}
func (r *Repo) CloneRepo() error {
	log.Println("Starting to clone repositories...")
	orgName, err := x.RepoURLParseOrgName(r.repoURL)
	if err != nil {
		return err
	}
	log.Printf("Parsed organization/subpath: %s", orgName)

	var repo []*platforms.RepoInfo

	log.Println("Fetching repositories from platform...")
	if orgName != "" {
		log.Printf("Organization/subpath detected: %s, fetching org repositories...", orgName)
		repo, err = r.platform.ListOrgRepo(r.ctx, orgName)
	} else {
		log.Println("No organization/subpath detected, fetching user repositories...")
		repo, err = r.platform.ListUserRepo(r.ctx)
	}
	if err != nil {
		return err
	}
	log.Printf("Successfully fetched %d repositories", len(repo))
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
				log.Printf("Failed to create Git instance for %s: %v", repo.CloneURL, err)
				return
			}
			if err := git.Clone(); err != nil {
				log.Printf("Failed to clone %s: %v", repo.CloneURL, err)
				return
			}
			log.Printf("Repository %s cloned successfully (took %s)", repo.CloneURL, time.Since(start))
		}(repo)
	}
	wg.Wait()
	return nil
}
func (r *Repo) RepoSync() error {
	log.Println("Starting repository sync...")
	targetURL, targetCredential, err := flags.GetTargetCredential(r.cmd)
	if err != nil {
		return err
	}
	log.Printf("Target URL parsed: %s", targetURL.String())
	targetPlatform, err := platforms.GetPlatform(targetURL, targetCredential)
	if err != nil {
		return err
	}
	log.Println("Target platform initialized successfully")

	orgName, err := x.RepoURLParseOrgName(r.repoURL)
	if err != nil {
		return err
	}
	var repos []*platforms.RepoInfo
	if orgName != "" {
		repos, err = r.platform.ListOrgRepo(r.ctx, orgName)
		log.Printf("Fetched %d repositories from org: %s", len(repos), orgName)
	} else {
		repos, err = r.platform.ListUserRepo(r.ctx)
		log.Printf("Fetched %d user repositories", len(repos))
	}
	if err != nil {
		return err
	}
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
			targetRepoURL := fmt.Sprintf("%s/%s.git", targetURL.String(), repo.Name)
			log.Printf("Source clone URL: %s", cloneURL)
			log.Printf("Target repo URL: %s", targetRepoURL)
			r.credential.CloneURL = cloneURL
			targetCredential.CloneURL = targetRepoURL
			doubleCredentialGit, err := NewDoubleCredentialGit(r.cmd, r.credential, targetCredential)
			if err != nil {
				log.Printf("Failed to create Git instance for %s: %v", repo.Name, err)
				return
			}
			log.Printf("Git instance created for %s", repo.Name)
			targetFullName, _ := targetCredential.GetFullName()
			detail, err := targetPlatform.GetRepoDetail(r.ctx, targetFullName)
			if err != nil || detail.ID == 0 {
				log.Printf("Target repository %s does not exist or cannot be fetched, creating...", repo.Name)
				if createErr := targetPlatform.CreateRepo(r.ctx, &platforms.RepoInfo{
					Name:        repo.Name,
					IsPrivate:   repo.IsPrivate,
					FullName:    targetFullName,
					Description: repo.Description,
					Homepage:    repo.Homepage,
				}); createErr != nil {
					log.Printf("Failed to create target repository %s: %v", repo.Name, createErr)
					return
				}
			}
			log.Printf("Pushing repository: %s", repo.Name)
			if err := doubleCredentialGit.Push(); err != nil {
				log.Printf("Failed to push %s: %v", repo.Name, err)
				return
			}
			log.Printf("Repository %s synced successfully (took %s)", repo.Name, time.Since(start))
		}(repo)
	}
	wg.Wait()
	log.Println("All repositories sync completed")
	return nil
}

func (r *Repo) CreateRelease() error {
	tags, err := r.platform.ListTags(r.ctx, r.fullName)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	for i, tag := range tags {
		start := time.Now()
		_, err := r.platform.CreateRelease(r.ctx, r.fullName, &platforms.ReleaseInfo{TagName: tag.TagName})
		if err != nil {
			log.Printf("Warning: failed to create release for tag '%s': %v", tag.TagName, err)
			continue
		}
		log.Printf("Created release for tag '%s' (%d/%d) in %s", tag.TagName, i+1, len(tags), time.Since(start))
	}
	return nil
}

func (r *Repo) Upload(tag string, filenames []string) error {
	log.Printf("Uploading files for tag '%s'...", tag)

	info, err := r.platform.GetTagReleaseInfo(r.ctx, r.fullName, tag)
	if err != nil {
		return fmt.Errorf("failed to get release info for tag '%s': %w", tag, err)
	}

	if err := r.platform.DeleteReleaseAssets(r.ctx, info, filenames); err != nil {
		log.Printf("Warning: failed to delete existing assets for tag '%s': %v", tag, err)
	}

	if err := r.platform.UploadReleaseAsset(r.ctx, info, filenames); err != nil {
		return fmt.Errorf("failed to upload assets for tag '%s': %w", tag, err)
	}

	log.Printf("Upload completed for tag '%s', %d files", tag, len(filenames))
	return nil
}

func (r *Repo) Download(tags []string) error {
	if len(tags) == 0 {
		return fmt.Errorf("no tags provided for download")
	}
	workspace, err := r.credential.GetCategoryNamWorkspace(credential.WorkspaceCategoryReleases, flags.GetWorkspace(r.cmd))
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}
	log.Printf("Downloading releases into workspace: %s", workspace)
	for i, tag := range tags {
		start := time.Now()
		log.Printf("Processing tag '%s' (%d/%d)", tag, i+1, len(tags))

		info, err := r.platform.GetTagReleaseInfo(r.ctx, r.fullName, tag)
		if err != nil {
			log.Printf("Warning: failed to get release info for tag '%s': %v", tag, err)
			continue
		}

		files, err := info.Download(workspace)
		if err != nil {
			log.Printf("Warning: failed to download files for tag '%s': %v", tag, err)
			continue
		}

		log.Printf("Downloaded %d files for tag '%s' in %s", len(files), tag, time.Since(start))
		for _, f := range files {
			log.Printf("  - %s", f)
		}
	}
	log.Printf("All downloads completed")
	return nil
}
