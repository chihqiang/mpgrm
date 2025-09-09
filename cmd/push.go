package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
	"time"
	"wangzhiqiang/mpgrm/factory"
	"wangzhiqiang/mpgrm/flags"
	"wangzhiqiang/mpgrm/pkg/logger"
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
				logger.Error("Failed to create Git instance: %v", err)
				return err
			}
			logger.Info("Starting push operation...")
			start := time.Now()
			if err := git.Push(); err != nil {
				logger.Error("Git push failed: %v", err)
				return err
			}
			elapsed := time.Since(start)
			logger.Info("Git push completed successfully, elapsed time: %s", elapsed)
			return nil
		},
	}
}
