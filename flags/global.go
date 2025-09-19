package flags

import (
	"github.com/urfave/cli/v3"
	"os"
	"path"
)

const (
	FlagsWorkspace = "workspace"
	FlagsPlatforms = "platform"
	FlagsEnvFile   = "env"

	FlagsFormUsername   = "username"
	FlagsFormPassword   = "password"
	FlagsFormToken      = "token"
	FlagsTargetUsername = "target-username"
	FlagsTargetPassword = "target-password"
	FlagsTargetToken    = "target-token"
)
const (
	EnvFormUserName = "USERNAME"
	EnvFormPassword = "PASSWORD"
	EnvFormToken    = "TOKEN"

	EnvTargetUserName = "TARGET_USERNAME"
	EnvTargetPassword = "TARGET_PASSWORD"
	EnvTargetToken    = "TARGET_TOKEN"
)

// GlobalFlags returns global flags applicable to all commands.
func GlobalFlags() []cli.Flag {
	dir, _ := os.Getwd()
	return []cli.Flag{
		&cli.StringFlag{
			Name:  FlagsWorkspace,
			Usage: "Where all the magic happens (local workspace for cache & temp files)",
			Value: path.Join(dir, "runtime"), // default workspace directory
		},
		&cli.StringFlag{
			Name:  FlagsPlatforms,
			Usage: "Path to the platform mapping file (for private deployments with custom domain mappings)",
			Value: path.Join(dir, "platforms.json"),
		},
		&cli.StringFlag{
			Name:  FlagsEnvFile,
			Usage: "Your secret sauce: path to the .env file",
			Value: ".env", // default env file
		},
		&cli.StringFlag{
			Name:    FlagsFormUsername,
			Usage:   "Username for authenticating with the source repository",
			Sources: cli.EnvVars(EnvFormUserName),
		},
		&cli.StringFlag{
			Name:    FlagsFormPassword,
			Usage:   "Password for authenticating with the source repository",
			Sources: cli.EnvVars(EnvFormPassword),
		},
		&cli.StringFlag{
			Name:    FlagsFormToken,
			Usage:   "Access token or secret for authenticating with the source repository",
			Sources: cli.EnvVars(EnvFormToken),
		},
		&cli.StringFlag{
			Name:    FlagsTargetUsername,
			Usage:   "Username for authenticating with the target repository",
			Sources: cli.EnvVars(EnvTargetUserName),
		},
		&cli.StringFlag{
			Name:    FlagsTargetPassword,
			Usage:   "Password for authenticating with the target repository",
			Sources: cli.EnvVars(EnvTargetPassword),
		},
		&cli.StringFlag{
			Name:    FlagsTargetToken,
			Usage:   "Access token or secret for authenticating with the target repository",
			Sources: cli.EnvVars(EnvTargetToken),
		},
	}
}

// GetWorkspace retrieves the workspace path from the command flags.
func GetWorkspace(cmd *cli.Command) string {
	return cmd.String(FlagsWorkspace)
}
