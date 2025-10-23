package cmd

import (
	"context"
	"github.com/chihqiang/mpgrm/factory"
	"github.com/chihqiang/mpgrm/flags"
	"github.com/chihqiang/mpgrm/pkg/logx"
	"github.com/urfave/cli/v3"
	"time"
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
				return err
			}
			logx.Info("Starting push operation...")
			start := time.Now()
			if err := git.Push(); err != nil {
				return err
			}
			elapsed := time.Since(start)
			logx.Info("Git push completed successfully, elapsed time: %s", elapsed)
			return nil
		},
	}
}
