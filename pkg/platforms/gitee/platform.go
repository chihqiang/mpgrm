package gitee

import (
	"fmt"
	"net/url"
	"wangzhiqiang/mpgrm/pkg/credential"
)

const (
	ApiURL = "https://gitee.com/api/v5/"
)

const (
	HOST        = "gitee.com"
	EnvPrefix   = "GITEE"
	EnvUsername = EnvPrefix + credential.EnvUsernameSuffix
	EnvToken    = EnvPrefix + credential.EnvTokenSuffix
)

type Platform struct {
	Credential *credential.Credential
	ApiURL     string
}

func (p *Platform) WithCredential(credential *credential.Credential) error {
	if credential == nil || credential.Token == "" {
		return fmt.Errorf("invalid %s credential: credential is nil or token is empty", EnvPrefix)
	}
	p.Credential = credential
	return nil
}

// GetURLWithToken 构建带 Token 和额外 query 参数的完整 API URL
func (p *Platform) GetURLWithToken(route string, query map[string]string) string {
	if p.ApiURL == "" {
		p.ApiURL = ApiURL
	}
	base, err := url.Parse(p.ApiURL)
	if err != nil {
		return ""
	}
	path, err := url.Parse(route)
	if err != nil {
		return ""
	}

	// 拼接 URL
	fullURL := base.ResolveReference(path)

	// 构建 query 参数
	q := fullURL.Query()
	if p.Credential != nil && p.Credential.Token != "" {
		q.Set("access_token", p.Credential.Token)
	}
	for k, v := range query {
		if v != "" {
			q.Set(k, v)
		}
	}
	fullURL.RawQuery = q.Encode()

	return fullURL.String()
}
