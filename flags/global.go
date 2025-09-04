package flags

import (
	"github.com/urfave/cli/v3"
	"os"
	"path"
)

const (
	FlagsWorkspace = "workspace"
	FlagsEnvFile   = "env"

	FlagsFormUsername   = "username"
	FlagsFormToken      = "token"
	FlagsTargetUsername = "target-username"
	FlagsTargetToken    = "target-token"
)
const (
	EnvFormUserName = "USERNAME"
	EnvFormToken    = "TOKEN"

	EnvTargetUserName = "TARGET_USERNAME"
	EnvTargetToken    = "TARGET_TOKEN"
)

func GlobalFlags() []cli.Flag {
	dir, _ := os.Getwd()
	return []cli.Flag{
		&cli.StringFlag{
			Name:  FlagsWorkspace,
			Usage: "Where all the magic happens (local workspace for cache & temp files)",
			Value: path.Join(dir, ".mpgrm"),
		},
		&cli.StringFlag{
			Name:  FlagsEnvFile,
			Usage: "Your secret sauce: path to the .env file",
			Value: ".env",
		},
		&cli.StringFlag{
			Name:    FlagsFormUsername,
			Usage:   "The wizard name for your source repo",
			Sources: cli.EnvVars(EnvFormUserName),
		},
		&cli.StringFlag{
			Name:    FlagsFormToken,
			Usage:   "Secret key or token to tame the source repo",
			Sources: cli.EnvVars(EnvFormToken),
		},
		&cli.StringFlag{
			Name:    FlagsTargetUsername,
			Usage:   "The hero name for your target repo",
			Sources: cli.EnvVars(EnvTargetUserName),
		},
		&cli.StringFlag{
			Name:    FlagsTargetToken,
			Usage:   "Secret key or token to rule the target repo",
			Sources: cli.EnvVars(EnvTargetToken),
		},
	}
}

func GetWorkspace(cmd *cli.Command) string {
	return cmd.String(FlagsWorkspace)
}
