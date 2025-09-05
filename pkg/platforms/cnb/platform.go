package cnb

import (
	"cnb.cool/cnb/sdk/go-cnb/cnb"
	"fmt"
	"wangzhiqiang/mpgrm/pkg/credential"
)

const (
	baseApiURL  = "https://api.cnb.cool"
	user        = "cnb"
	downloadURL = "https://cnb.cool"
	cloneURL    = "https://cnb.cool"
)

const (
	HOST = "cnb.cool"
)

const (
	EnvPrefix   = "CNB"
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

func (p *Platform) GetClient() (*cnb.Client, error) {
	if p.Credential.ApiURL == "" {
		p.Credential.ApiURL = baseApiURL
	}
	return cnb.NewClient(nil).WithAuthToken(p.Credential.Token).WithURLs(p.Credential.ApiURL)
}
