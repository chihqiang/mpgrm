package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"time"
	"wangzhiqiang/mpgrm/factory"
	"wangzhiqiang/mpgrm/flags"
)

func ReleasesCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "releases",
		Usage:                  "Manage your releases with ease",
		Flags:                  flags.FormTargetRelease(),
		Commands: []*cli.Command{
			{
				Name:  "upload",
				Usage: "Upload the chosen release version",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					log.Println("Starting upload...")
					repo, err := factory.NewRepo(ctx, cmd)
					if err != nil {
						return fmt.Errorf("failed to initialize repo: %w", err)
					}

					tag, err := flags.GetFirstTags(cmd)
					if err != nil {
						return fmt.Errorf("failed to get tag: %w", err)
					}

					files := flags.GetFiles(cmd)
					log.Printf("Uploading %d files for tag '%s'...", len(files), tag)
					if err := repo.Upload(tag, files); err != nil {
						return fmt.Errorf("upload failed for tag '%s': %w", tag, err)
					}

					log.Printf("Upload completed for tag '%s' in %s", tag, time.Since(start))
					return nil
				},
			},
			{
				Name:  "download",
				Usage: "Download release attachments for specified tags",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					log.Println("Starting download...")
					repo, err := factory.NewRepo(ctx, cmd)
					if err != nil {
						return fmt.Errorf("failed to initialize repo: %w", err)
					}

					tags := flags.GetTags(cmd)
					log.Printf("Downloading releases for %d tag(s)...", len(tags))
					if err := repo.Download(tags); err != nil {
						return fmt.Errorf("download failed: %w", err)
					}

					log.Printf("Download completed in %s", time.Since(start))
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "Create releases for all tags in the repo",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					log.Println("Starting release creation...")
					repo, err := factory.NewRepo(ctx, cmd)
					if err != nil {
						return fmt.Errorf("failed to initialize repo: %w", err)
					}

					if err := repo.CreateRelease(); err != nil {
						return fmt.Errorf("release creation failed: %w", err)
					}
					log.Printf("All releases created successfully in %s", time.Since(start))
					return nil
				},
			},
			{
				Name:  "sync",
				Usage: "Sync releases from source repo to target repo",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					log.Println("Starting release sync...")
					target, err := factory.NewDoubleRepo(ctx, cmd)
					if err != nil {
						return fmt.Errorf("failed to initialize target repo: %w", err)
					}

					tags := flags.GetTags(cmd)
					log.Printf("Syncing releases for %d tag(s)...", len(tags))
					if err := target.ReleaseSync(tags); err != nil {
						return fmt.Errorf("release sync failed: %w", err)
					}

					log.Printf("Release sync completed in %s", time.Since(start))
					return nil
				},
			},
		},
	}
}
