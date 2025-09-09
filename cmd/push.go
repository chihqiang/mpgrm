package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
	"log"
	"time"
	"wangzhiqiang/mpgrm/factory"
	"wangzhiqiang/mpgrm/flags"
)

// PushCommand defines the CLI command to push repositories to the target.
func PushCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "push",
		Usage:                  "Stress-free push, releases on demand",
		Flags:                  flags.FormTargetRepoPush(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			git, err := factory.NewCmdDoubleGit(ctx, cmd)
			if err != nil {
				log.Printf("Failed to create Git instance: %v", err)
				return err
			}
			log.Println("Starting push operation...")
			start := time.Now()
			if err := git.Push(); err != nil {
				log.Printf("Git push failed: %v", err)
				return err
			}
			elapsed := time.Since(start)
			log.Printf("Git push completed successfully, elapsed time: %s", elapsed)
			return nil
		},
	}
}
