package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
	"log"
	"time"
	"wangzhiqiang/mpgrm/factory"
	"wangzhiqiang/mpgrm/flags"
)

// PushCommand 返回一个 CLI 命令，用于执行 Git 仓库迁移操作。
// 它包括以下步骤：
// 1. 根据提供的 form 仓库和目标仓库的凭据初始化 GitMigrate 实例。
// 2. 生成本地 workspace 用于存放克隆的仓库。
// 3. 克隆指定的分支和标签到本地 workspace。
// 4. 将本地仓库推送到目标仓库。
func PushCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "push",
		Usage:                  "Stress-free push, releases on demand",
		Flags:                  flags.FormTargetRepo(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			git, err := factory.NewGitCmd(ctx, cmd)
			if err != nil {
				log.Printf("Failed to create Git instance: %v", err)
				return err
			}
			log.Println("Git instance created successfully")

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
