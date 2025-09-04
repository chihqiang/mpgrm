package gitea

import (
	"code.gitea.io/sdk/gitea"
	"fmt"
	"wangzhiqiang/mpgrm/pkg/credential"
)

const (
	ApiURL      = "https://gitea.com"
	HOST        = "gitea.com"
	EnvPrefix   = "GITEA"
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
func (p *Platform) GetClient() (*gitea.Client, error) {
	if p.ApiURL == "" {
		p.ApiURL = ApiURL
	}
	return gitea.NewClient(p.ApiURL, gitea.SetToken(p.Credential.Username))
}
