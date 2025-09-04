package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"runtime"
	"wangzhiqiang/mpgrm/cmd"
	"wangzhiqiang/mpgrm/flags"
	"wangzhiqiang/mpgrm/pkg/platforms"
	"wangzhiqiang/mpgrm/pkg/platforms/cnb"
	"wangzhiqiang/mpgrm/pkg/platforms/gitea"
	"wangzhiqiang/mpgrm/pkg/platforms/gitee"
	"wangzhiqiang/mpgrm/pkg/platforms/github"
)

var (
	version  = "main"
	commands []*cli.Command
)

func init() {
	platforms.Register("cnb.cool", cnb.EnvPrefix, func() platforms.IPlatform {
		return &cnb.Platform{}
	})
	platforms.Register("gitea.com", gitea.EnvPrefix, func() platforms.IPlatform {
		return &gitea.Platform{}
	})
	platforms.Register("gitee.com", gitee.EnvPrefix, func() platforms.IPlatform {
		return &gitee.Platform{}
	})
	platforms.Register("github.com", github.EnvPrefix, func() platforms.IPlatform {
		return &gitee.Platform{}
	})
	commands = append(commands, cmd.PushCommand())
	commands = append(commands, cmd.ReleasesCommand())
	commands = append(commands, cmd.RepoCommand())
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
		return ctx, nil
	}
	app.Commands = commands
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
