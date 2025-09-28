package cmd

import (
	"context"
	"github.com/chihqiang/mpgrm/flags"
	"github.com/chihqiang/mpgrm/pkg/logger"
	"github.com/urfave/cli/v3"
)

func CredentialCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "credential",
		Usage:                  "Manage repository credentials with ease and push anytime",

		Commands: []*cli.Command{
			{
				Name:  "show",
				Usage: "Show repository credentials",
				Flags: flags.CredentialFlags(),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					url, cred, err := flags.GetFormCredential(cmd, flags.UseEnv(cmd))
					if err != nil {
						return err
					}
					logger.Info("Repository: %s", url)
					logger.Info("Credential: %s", cred.String())
					return nil
				},
			},
		},
	}
}
