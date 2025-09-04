package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
	"log"
	"time"
	"wangzhiqiang/mpgrm/factory"
	"wangzhiqiang/mpgrm/flags"
)

func RepoCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "repo",
		Usage:                  "Your all-in-one repository command center",
		Flags:                  flags.FormTargetRepo(),
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all available repositories",
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
				Name:  "sync",
				Usage: "Keep your org or personal repos marching in step",
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
