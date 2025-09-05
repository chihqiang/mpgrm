package gitea

import (
	"code.gitea.io/sdk/gitea"
	"fmt"
	"wangzhiqiang/mpgrm/pkg/credential"
)

const (
	HOST        = "gitea.com"
	EnvPrefix   = "GITEA"
	EnvUsername = EnvPrefix + credential.EnvUsernameSuffix
	EnvToken    = EnvPrefix + credential.EnvTokenSuffix
)

type Platform struct {
	Credential *credential.Credential
}

func (p *Platform) WithCredential(credential *credential.Credential) error {
	if credential == nil || credential.Token == "" {
		return fmt.Errorf("invalid %s credential: credential is nil or token is empty", EnvPrefix)
	}
	p.Credential = credential
	return nil
}
func (p *Platform) GetClient() (*gitea.Client, error) {
	if p.Credential.ApiURL == "" {
		p.Credential.ApiURL = "https://gitea.com"
	}
	return gitea.NewClient(p.Credential.ApiURL, gitea.SetToken(p.Credential.Username))
}
