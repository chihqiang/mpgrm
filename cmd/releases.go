package cmd

import (
	"context"
	"fmt"
	"github.com/chihqiang/mpgrm/factory"
	"github.com/chihqiang/mpgrm/flags"
	"github.com/chihqiang/mpgrm/pkg/logger"
	"github.com/urfave/cli/v3"
	"time"
)

// ReleasesCommand defines the CLI command to manage releases.
func ReleasesCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "releases",
		Usage:                  "Manage your releases with ease",
		Commands: []*cli.Command{
			{
				Name:  "upload",
				Usage: "Upload the chosen release version",
				Flags: flags.FormReleaseUploadFiles(),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					logger.Info("Starting upload...")
					repo, err := factory.NewRepo(ctx, cmd)
					if err != nil {
						return fmt.Errorf("failed to initialize repo: %w", err)
					}

					tag, err := flags.GetFirstTags(cmd)
					if err != nil {
						return fmt.Errorf("failed to get tag: %w", err)
					}

					files := flags.GetFiles(cmd)
					logger.Info("Uploading %d files for tag '%s'...", len(files), tag)
					if err := repo.Upload(tag, files); err != nil {
						return fmt.Errorf("upload failed for tag '%s': %w", tag, err)
					}

					logger.Info("Upload completed for tag '%s' in %s", tag, time.Since(start))
					return nil
				},
			},
			{
				Name:  "download",
				Usage: "Download release attachments for specified tags",
				Flags: flags.FormReleaseDownload(),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					logger.Info("Starting download...")
					repo, err := factory.NewRepo(ctx, cmd)
					if err != nil {
						return fmt.Errorf("failed to initialize repo: %w", err)
					}

					tags := flags.GetTags(cmd)
					logger.Info("Downloading releases for %d tag(s)...", len(tags))
					if _, err := repo.Download(tags); err != nil {
						return fmt.Errorf("download failed: %w", err)
					}

					logger.Info("Download completed in %s", time.Since(start))
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "Create releases for all tags in the repo",
				Flags: flags.FormFlags(),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					logger.Info("Starting release creation...")
					repo, err := factory.NewRepo(ctx, cmd)
					if err != nil {
						return fmt.Errorf("failed to initialize repo: %w", err)
					}

					if err := repo.CreateRelease(); err != nil {
						return fmt.Errorf("release creation failed: %w", err)
					}
					logger.Info("All releases created successfully in %s", time.Since(start))
					return nil
				},
			},
			{
				Name:  "sync",
				Flags: flags.FormTargetReleaseSync(),
				Usage: "Sync releases from source repo to target repo",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					logger.Info("Starting release sync...")
					target, err := factory.NewDoubleRepo(ctx, cmd)
					if err != nil {
						return fmt.Errorf("failed to initialize target repo: %w", err)
					}

					tags := flags.GetTags(cmd)
					logger.Info("Syncing releases for %d tag(s)...", len(tags))
					if err := target.ReleaseSync(tags); err != nil {
						return fmt.Errorf("release sync failed: %w", err)
					}
					logger.Info("Release sync completed in %s", time.Since(start))
					return nil
				},
			},
		},
	}
}
