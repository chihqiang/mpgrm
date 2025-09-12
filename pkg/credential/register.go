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
	EnvPasswordSuffix = "_PASSWORD"
	EnvTokenSuffix    = "_TOKEN"
)

func Register(k, v string) {
	hostEnvPrefix[k] = v
}

func GetCredential(repo *url.URL, username, password, token string, readEnv bool) (*Credential, error) {
	c := &Credential{CloneURL: repo.String()}
	if username != "" && (password != "" || token != "") {
		c.Username = username
		c.Password = password
		c.Token = token
		return c, nil
	}
	// 2. 尝试从 repo.User 提取
	if repo.User != nil {
		if u := repo.User.Username(); u != "" {
			username = u
		}
		if p, ok := repo.User.Password(); ok {
			password = p
		}
	}
	if username != "" && password != "" {
		c.Username = username
		c.Password = password
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
		c.Username = os.Getenv(envPrefix + EnvUsernameSuffix)
		c.Password = os.Getenv(envPrefix + EnvPasswordSuffix)
		c.Token = os.Getenv(envPrefix + EnvTokenSuffix)
		if c.Username != "" && (c.Password != "" || c.Token != "") {
			return c, nil
		}
	}
	return c, fmt.Errorf("missing credential: username/token not found for repo %q", repo.String())
}
