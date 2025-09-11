package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"os"
	"runtime"
	"wangzhiqiang/mpgrm/cmd"
	"wangzhiqiang/mpgrm/flags"
	"wangzhiqiang/mpgrm/pkg/logger"
	"wangzhiqiang/mpgrm/register"
)

var (
	version  = "main"
	commands []*cli.Command
)

func init() {
	commands = append(commands, cmd.PushCommand())
	commands = append(commands, cmd.ReleasesCommand())
	commands = append(commands, cmd.RepoCommand())
	commands = append(commands, cmd.CredentialCommand())
}

func main() {
	app := &cli.Command{}
	app.Name = "mpgrm"
	app.Usage = "own your Git repos, effortlessly across platforms"
	app.Version = version
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("mpgrm version %s %s/%s\n", cmd.Version, runtime.GOOS, runtime.GOARCH)
	}
	app.Flags = flags.GlobalFlags()
	app.Before = func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		_ = godotenv.Load(cmd.String(flags.FlagsEnvFile))
		register.Platforms(cmd)
		return ctx, nil
	}
	app.Commands = commands
	if err := app.Run(context.Background(), os.Args); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
