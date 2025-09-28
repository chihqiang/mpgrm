package register

import (
	"encoding/json"
	"github.com/chihqiang/mpgrm/flags"
	"github.com/chihqiang/mpgrm/pkg/platforms"
	"github.com/chihqiang/mpgrm/pkg/platforms/cnb"
	"github.com/chihqiang/mpgrm/pkg/platforms/gitea"
	"github.com/chihqiang/mpgrm/pkg/platforms/gitee"
	"github.com/chihqiang/mpgrm/pkg/platforms/github"
	"github.com/urfave/cli/v3"
	"net/url"
	"os"
	"slices"
)

type Config struct {
	Host   string `json:"host"`
	ApiURL string `json:"api_url"`
}

// PFactory
// 平台注册信息
type PFactory struct {
	EnvPrefix string
	Factory   func(apiURL string) platforms.IPlatform
}

// 所有支持的平台配置
var pFactories = map[string]PFactory{
	gitee.HOST: {
		EnvPrefix: gitee.EnvPrefix,
		Factory:   func(apiURL string) platforms.IPlatform { return &gitee.Platform{} },
	},
	github.HOST: {
		EnvPrefix: github.EnvPrefix,
		Factory:   func(apiURL string) platforms.IPlatform { return &github.Platform{} },
	},
	//私有化
	cnb.HOST: {
		EnvPrefix: cnb.EnvPrefix,
		Factory:   func(apiURL string) platforms.IPlatform { return &cnb.Platform{ApiURL: apiURL} },
	},
	gitea.HOST: {
		EnvPrefix: gitea.EnvPrefix,
		Factory:   func(apiURL string) platforms.IPlatform { return &gitea.Platform{ApiURL: apiURL} },
	},
}

var (
	allowPrivateHosts = []string{cnb.HOST, gitea.HOST}
)

func Platforms(cmd *cli.Command) {
	// 注册默认平台
	for host, cfg := range pFactories {
		registerPlatform(host, "", cfg)
	}
	// 尝试加载自定义配置
	if data, err := os.ReadFile(cmd.String(flags.FlagsPlatforms)); err == nil {
		var mapping []Config
		if err := json.Unmarshal(data, &mapping); err == nil {
			for _, hp := range mapping {
				apiURL, err := url.Parse(hp.ApiURL)
				if err != nil {
					continue
				}
				//需要在允许私有化部署的host
				if !slices.Contains(allowPrivateHosts, hp.Host) {
					continue
				}
				if cfg, ok := pFactories[hp.Host]; ok {
					registerPlatform(apiURL.Host, apiURL.String(), cfg)
				}
			}
		}
	}
}

// 注册逻辑
func registerPlatform(host, apiURL string, pf PFactory) {
	platforms.Register(host, pf.EnvPrefix, func() platforms.IPlatform {
		return pf.Factory(apiURL)
	})
}
