package credential

import (
	"fmt"
	"net/url"
	"os"
)

var (
	hostEnvPrefix = map[string]string{}
)

const (
	EnvUsernameSuffix = "_USERNAME"
	EnvTokenSuffix    = "_TOKEN"
)

func Register(k, v string) {
	hostEnvPrefix[k] = v
}

func GetCredential(repo *url.URL, username, token string, readEnv bool) (*Credential, error) {
	c := &Credential{CloneURL: repo.String()}
	// 1. 优先使用传入参数
	if username != "" && token != "" {
		c.Username = username
		c.Token = token
		return c, nil
	}

	// 2. 尝试从 repo.User 提取
	if repo.User != nil {
		if username == "" {
			username = repo.User.Username()
		}
		if token == "" {
			if p, ok := repo.User.Password(); ok {
				token = p
			}
		}
	}
	if username != "" && token != "" {
		c.Username = username
		c.Token = token
		return c, nil
	}

	// 3. 根据 host 查找环境变量（可选）
	if readEnv {
		host := repo.Host
		if host == "" {
			return c, fmt.Errorf("missing credential: unknown host (repo=%q)", repo.Redacted())
		}
		envPrefix, ok := hostEnvPrefix[host]
		if !ok {
			return nil, fmt.Errorf("missing credential: unsupported platform %q (host=%q)", host, host)
		}

		if username == "" {
			c.Username = os.Getenv(envPrefix + EnvUsernameSuffix)
		}
		if token == "" {
			c.Token = os.Getenv(envPrefix + EnvTokenSuffix)
		}
		if c.Username != "" && c.Token != "" {
			return c, nil
		}
	}
	return c, fmt.Errorf("missing credential: username/token not found for repo %q", repo.String())
}
