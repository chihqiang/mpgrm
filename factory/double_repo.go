package factory

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"time"
	"wangzhiqiang/mpgrm/flags"
	"wangzhiqiang/mpgrm/pkg/credential"
	"wangzhiqiang/mpgrm/pkg/logger"
	"wangzhiqiang/mpgrm/pkg/platforms"
)

// DoubleRepo represents a repository sync context between a source and a target repository.
// It contains context, CLI command, source and target platforms, credentials, and repository full names.
type DoubleRepo struct {
	ctx context.Context // Context for controlling cancellation and deadlines
	cmd *cli.Command    // CLI command instance

	platform   platforms.IPlatform    // Source repository platform interface (GitHub, Gitee, Gitea, etc.)
	credential *credential.Credential // Source repository authentication credential
	fullName   string                 // Full name of the source repository (owner/repo[/subdir])

	targetPlatform   platforms.IPlatform    // Target repository platform interface
	targetCredential *credential.Credential // Target repository authentication credential
	targetFullName   string                 // Full name of the target repository (owner/repo[/subdir])
}

// NewDoubleRepo initializes a DoubleRepo instance with source and target repository information.
// It sets up the context, CLI command, credentials, platform interfaces, and full repository names.
func NewDoubleRepo(ctx context.Context, cmd *cli.Command) (*DoubleRepo, error) {
	rt := &DoubleRepo{ctx: ctx, cmd: cmd}

	// Get source repository URL and credential
	repoURL, cred, err := flags.GetFormCredential(cmd, true)
	if err != nil {
		return rt, err
	}
	rt.credential = cred

	// Get source platform interface
	platform, err := platforms.GetPlatform(repoURL, cred)
	if err != nil {
		return rt, err
	}
	rt.platform = platform
	// Parse source repository full name
	fullName, err := cred.GetFullName()
	if err != nil {
		return rt, err
	}
	rt.fullName = fullName

	// Get target repository URL and credential
	targetRepoURL, targetCred, err := flags.GetTargetCredential(cmd)
	if err != nil {
		return rt, err
	}
	rt.targetCredential = targetCred

	// Parse target repository full name
	targetFullName, err := targetCred.GetFullName()
	if err != nil {
		return rt, err
	}
	rt.targetFullName = targetFullName

	// Get target platform interface
	targetPlatform, err := platforms.GetPlatform(targetRepoURL, targetCred)
	if err != nil {
		return rt, err
	}
	rt.targetPlatform = targetPlatform

	return rt, nil
}

// ReleaseSync synchronizes releases from the source repository to the target repository.
// It can optionally filter by specific tags provided in the `tags` slice.
// Parameters:
//   - tags: a slice of tag names to be synchronized. If empty, all tags will be considered.
//
// Returns:
//   - error: any error encountered during the synchronization process.
func (t *DoubleRepo) ReleaseSync(tags []string) error {
	if len(tags) == 0 {
		return fmt.Errorf("no tags provided for release sync")
	}

	workspace, err := t.credential.GetCategoryNamWorkspace(credential.WorkspaceCategoryReleases, flags.GetWorkspace(t.cmd))
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	for i, tag := range tags {
		start := time.Now()
		logger.Info("Processing tag %s (%d/%d)", tag, i+1, len(tags))

		// 获取源 Release 信息
		info, err := t.platform.GetTagReleaseInfo(t.ctx, t.fullName, tag)
		if err != nil {
			logger.Warning("failed to get source release info for tag '%s': %v", tag, err)
			continue
		}

		// 下载源 Release 文件
		files, err := info.Download(workspace)
		if err != nil {
			logger.Warning("download release files for tag '%s': %v", tag, err)
			continue
		}
		logger.Info("Downloaded %d files for tag '%s' in %s", len(files), tag, time.Since(start))

		// 获取目标 Release，如果不存在则创建
		releaseInfo, err := t.targetPlatform.GetTagReleaseInfo(t.ctx, t.targetFullName, tag)
		if err != nil {
			releaseInfo, err = t.targetPlatform.CreateRelease(t.ctx, t.targetFullName, &platforms.ReleaseInfo{TagName: tag})
			if err != nil {
				logger.Warning("failed to create target release for tag '%s': %v", tag, err)
				continue
			}
			logger.Info("Created target release for tag '%s'", tag)
		}

		// 上传文件到目标 Release
		if err := t.targetPlatform.UploadReleaseAsset(t.ctx, releaseInfo, files); err != nil {
			logger.Warning("failed to upload files to target release for tag '%s': %v", tag, err)
			continue
		}
		logger.Info("Uploaded %d files to target release for tag '%s'", len(files), tag)
	}
	logger.Info("Release sync completed for %d tags", len(tags))
	return nil
}
