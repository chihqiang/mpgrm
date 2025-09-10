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

	targetPlatform   platforms.IPlatform    // Target repository platform interface
	targetCredential *credential.Credential // Target repository authentication credential
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
	// Get target repository URL and credential
	targetRepoURL, targetCred, err := flags.GetTargetCredential(cmd)
	if err != nil {
		return rt, err
	}
	rt.targetCredential = targetCred

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

	start := time.Now()
	var successCount, failCount int
	var failedTags []string

	targetFullName, err := t.targetCredential.GetFullName()
	if err != nil {
		return fmt.Errorf("failed to get target full name: %w", err)
	}

	for i, tag := range tags {
		tagStart := time.Now()
		logger.Info("Processing tag %s (%d/%d)", tag, i+1, len(tags))

		// 初始化 repo
		repo, err := NewRepo(t.ctx, t.cmd)
		if err != nil {
			logger.Warning("failed to init repo for tag '%s': %v", tag, err)
			failCount++
			failedTags = append(failedTags, tag)
			logger.Info("Tag %s completed with errors, elapsed: %s", tag, time.Since(tagStart))
			continue
		}

		mapFiles, err := repo.Download([]string{tag})
		if err != nil {
			logger.Warning("failed to download files for tag '%s': %v", tag, err)
			failCount++
			failedTags = append(failedTags, tag)
			logger.Info("Tag %s completed with errors, elapsed: %s", tag, time.Since(tagStart))
			continue
		}

		files, ok := mapFiles[tag]
		if !ok || len(files) == 0 {
			logger.Warning("no files found for tag '%s'", tag)
			failCount++
			failedTags = append(failedTags, tag)
			logger.Info("Tag %s completed with errors, elapsed: %s", tag, time.Since(tagStart))
			continue
		}

		// 获取目标 Release
		releaseInfo, err := t.targetPlatform.GetTagReleaseInfo(t.ctx, targetFullName, tag)
		if err != nil {
			logger.Info("Release for tag '%s' not found, creating...", tag)
			releaseInfo, err = t.targetPlatform.CreateRelease(t.ctx, targetFullName, &platforms.ReleaseInfo{
				TagName: tag,
			})
			if err != nil {
				logger.Warning("failed to create target release for tag '%s': %v", tag, err)
				failCount++
				failedTags = append(failedTags, tag)
				logger.Info("Tag %s completed with errors, elapsed: %s", tag, time.Since(tagStart))
				continue
			}
			logger.Info("Created target release for tag '%s'", tag)
		}

		// 上传文件
		if err := t.targetPlatform.UploadReleaseAsset(t.ctx, releaseInfo, files); err != nil {
			logger.Warning("failed to upload %d files to release for tag '%s': %v", len(files), tag, err)
			failCount++
			failedTags = append(failedTags, tag)
			logger.Info("Tag %s completed with errors, elapsed: %s", tag, time.Since(tagStart))
			continue
		}

		// 成功日志里带耗时
		successCount++
		logger.Info("Uploaded %d files to release for tag '%s', elapsed: %s", len(files), tag, time.Since(tagStart))
	}

	elapsed := time.Since(start)
	if failCount > 0 {
		logger.Warning("Release sync completed: %d success, %d failed (%v), total %d tags, total elapsed: %s", successCount, failCount, failedTags, len(tags), elapsed)
	} else {
		logger.Info("Release sync completed: %d success, %d failed, total %d tags, total elapsed: %s", successCount, failCount, len(tags), elapsed)
	}
	return nil
}
