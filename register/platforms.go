package register

import (
	"encoding/json"
	"github.com/urfave/cli/v3"
	"net/url"
	"os"
	"wangzhiqiang/mpgrm/flags"
	"wangzhiqiang/mpgrm/pkg/platforms"
	"wangzhiqiang/mpgrm/pkg/platforms/cnb"
	"wangzhiqiang/mpgrm/pkg/platforms/gitea"
	"wangzhiqiang/mpgrm/pkg/platforms/gitee"
	"wangzhiqiang/mpgrm/pkg/platforms/github"
)

type HostPlatform struct {
	PlatformHost string `json:"platform_host"`
	ApiURL       string `json:"api_url"`
}

// 平台注册信息
type platformConfig struct {
	EnvPrefix string
	Factory   func(apiURL string) platforms.IPlatform
}

// 所有支持的平台配置
var platformConfigs = map[string]platformConfig{
	cnb.HOST: {
		EnvPrefix: cnb.EnvPrefix,
		Factory:   func(apiURL string) platforms.IPlatform { return &cnb.Platform{ApiURL: apiURL} },
	},
	gitea.HOST: {
		EnvPrefix: gitea.EnvPrefix,
		Factory:   func(apiURL string) platforms.IPlatform { return &gitea.Platform{ApiURL: apiURL} },
	},
	gitee.HOST: {
		EnvPrefix: gitee.EnvPrefix,
		Factory:   func(apiURL string) platforms.IPlatform { return &gitee.Platform{ApiURL: apiURL} },
	},
	github.HOST: {
		EnvPrefix: github.EnvPrefix,
		Factory:   func(apiURL string) platforms.IPlatform { return &github.Platform{ApiURL: apiURL} },
	},
}

func Platforms(cmd *cli.Command) {
	// 注册默认平台
	for host, cfg := range platformConfigs {
		registerPlatform(host, "", cfg)
	}

	// 尝试加载自定义配置
	if data, err := os.ReadFile(cmd.String(flags.FlagsPlatforms)); err == nil {
		var mapping []HostPlatform
		if err := json.Unmarshal(data, &mapping); err == nil {
			for _, hp := range mapping {
				apiURL, err := url.Parse(hp.ApiURL)
				if err != nil {
					continue
				}
				if cfg, ok := platformConfigs[hp.PlatformHost]; ok {
					registerPlatform(apiURL.Host, apiURL.String(), cfg)
				}
			}
		}
	}
}

// 注册逻辑
func registerPlatform(host, apiURL string, cfg platformConfig) {
	platforms.Register(host, cfg.EnvPrefix, func() platforms.IPlatform {
		return cfg.Factory(apiURL)
	})
}
