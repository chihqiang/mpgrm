package factory

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"time"
	flags2 "wangzhiqiang/mpgrm/flags"
	"wangzhiqiang/mpgrm/pkg/credential"
	"wangzhiqiang/mpgrm/pkg/platforms"
	"wangzhiqiang/mpgrm/pkg/x"
)

type DoubleRepo struct {
	ctx context.Context
	cmd *cli.Command

	platform         platforms.IPlatform
	credential       *credential.Credential
	fullName         string
	targetPlatform   platforms.IPlatform
	targetCredential *credential.Credential
	targetFullName   string
}

func NewDoubleRepo(ctx context.Context, cmd *cli.Command) (*DoubleRepo, error) {
	rt := &DoubleRepo{ctx: ctx, cmd: cmd}
	repoURL, cred, err := flags2.GetFormCredential(cmd, true)
	if err != nil {
		return rt, err
	}
	rt.credential = cred

	platform, err := platforms.GetPlatform(repoURL, cred)
	if err != nil {
		return rt, err
	}
	rt.platform = platform

	fullName, err := x.RepoURLParseFullName(repoURL)
	if err != nil {
		return rt, err
	}
	rt.fullName = fullName

	targetRepoURL, targetCred, err := flags2.GetTargetCredential(cmd)
	if err != nil {
		return rt, err
	}
	rt.targetCredential = targetCred
	targetFullName, err := x.RepoURLParseFullName(targetRepoURL)
	if err != nil {
		return rt, err
	}
	rt.targetFullName = targetFullName

	targetPlatform, err := platforms.GetPlatform(targetRepoURL, targetCred)
	if err != nil {
		return rt, err
	}
	rt.targetPlatform = targetPlatform

	return rt, nil
}
func (t *DoubleRepo) ReleaseSync(tags []string) error {
	if len(tags) == 0 {
		return fmt.Errorf("no tags provided for release sync")
	}

	workspace, err := t.credential.GetCategoryNamWorkspace(credential.WorkspaceCategoryReleases, flags2.GetWorkspace(t.cmd))
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	for i, tag := range tags {
		start := time.Now()
		log.Printf("Processing tag %s (%d/%d)", tag, i+1, len(tags))

		// 获取源 Release 信息
		info, err := t.platform.GetTagReleaseInfo(t.ctx, t.fullName, tag)
		if err != nil {
			log.Printf("Warning: failed to get source release info for tag '%s': %v", tag, err)
			continue
		}

		// 下载源 Release 文件
		files, err := info.Download(workspace)
		if err != nil {
			log.Printf("Warning: failed to download release files for tag '%s': %v", tag, err)
			continue
		}
		log.Printf("Downloaded %d files for tag '%s' in %s", len(files), tag, time.Since(start))

		// 获取目标 Release，如果不存在则创建
		releaseInfo, err := t.targetPlatform.GetTagReleaseInfo(t.ctx, t.targetFullName, tag)
		if err != nil {
			releaseInfo, err = t.targetPlatform.CreateRelease(t.ctx, t.targetFullName, &platforms.ReleaseInfo{TagName: tag})
			if err != nil {
				log.Printf("Warning: failed to create target release for tag '%s': %v", tag, err)
				continue
			}
			log.Printf("Created target release for tag '%s'", tag)
		}

		// 上传文件到目标 Release
		if err := t.targetPlatform.UploadReleaseAsset(t.ctx, releaseInfo, files); err != nil {
			log.Printf("Warning: failed to upload files to target release for tag '%s': %v", tag, err)
			continue
		}
		log.Printf("Uploaded %d files to target release for tag '%s'", len(files), tag)
	}

	log.Printf("Release sync completed for %d tags", len(tags))
	return nil
}
