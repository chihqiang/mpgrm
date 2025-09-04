package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
	"log"
	"time"
	"wangzhiqiang/mpgrm/factory"
	"wangzhiqiang/mpgrm/flags"
)

// RepoCommand defines the CLI command for repository operations.
func RepoCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "repo",
		Usage:                  "Your all-in-one repository command center",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all available repositories",
				Flags: flags.FormFlags(),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					repo, err := factory.NewRepo(ctx, cmd)
					if err != nil {
						return err
					}
					log.Println("Repo initialized successfully")
					if err := repo.ListRepo(); err != nil {
						return err
					}
					log.Printf("Repository listing completed successfully in %s", time.Since(start))
					return nil
				},
			},
			{
				Name:  "clone",
				Usage: "Pull down the code and make it yours",
				Flags: flags.FormFlags(),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					repo, err := factory.NewRepo(ctx, cmd)
					if err != nil {
						return err
					}
					if err := repo.CloneRepo(); err != nil {
						return err
					}
					elapsed := time.Since(start)
					log.Printf("clone completed successfully in %s", elapsed)
					return nil
				},
			},
			{
				Name:  "sync",
				Usage: "Keep your org or personal repos marching in step",
				Flags: flags.FormTargetRepo(),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					start := time.Now()
					repo, err := factory.NewRepo(ctx, cmd)
					if err != nil {
						return err
					}
					log.Println("Repo initialized successfully")
					if err := repo.RepoSync(); err != nil {
						return err
					}
					log.Printf("RepoSync completed successfully (elapsed: %s)", time.Since(start))
					return nil
				},
			},
		},
	}
}
