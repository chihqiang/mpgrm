package register

import (
	"encoding/json"
	"github.com/urfave/cli/v3"
	"os"
	"wangzhiqiang/mpgrm/flags"
	"wangzhiqiang/mpgrm/pkg/platforms"
	"wangzhiqiang/mpgrm/pkg/platforms/cnb"
	"wangzhiqiang/mpgrm/pkg/platforms/gitea"
	"wangzhiqiang/mpgrm/pkg/platforms/gitee"
	"wangzhiqiang/mpgrm/pkg/platforms/github"
)

func Platforms(cmd *cli.Command) {
	// 官方域名 -> 注册逻辑
	defaults := map[string]func(string){
		cnb.HOST: func(d string) {
			platforms.Register(d, cnb.EnvPrefix, func() platforms.IPlatform { return &cnb.Platform{} })
		},
		gitea.HOST: func(d string) {
			platforms.Register(d, gitea.EnvPrefix, func() platforms.IPlatform { return &gitea.Platform{} })
		},
		gitee.HOST: func(d string) {
			platforms.Register(d, gitee.EnvPrefix, func() platforms.IPlatform { return &gitee.Platform{} })
		},
		github.HOST: func(d string) {
			platforms.Register(d, github.EnvPrefix, func() platforms.IPlatform { return &github.Platform{} })
		},
	}

	// 注册默认域名
	for domain, fn := range defaults {
		fn(domain)
	}

	// 尝试加载自定义配置
	if data, err := os.ReadFile(cmd.String(flags.FlagsPlatforms)); err == nil {
		var mapping map[string]string
		if json.Unmarshal(data, &mapping) == nil {
			for custom, official := range mapping {
				if fn, ok := defaults[official]; ok {
					fn(custom)
				}
			}
		}
	}
}
