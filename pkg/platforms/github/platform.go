package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v73/github"
	"golang.org/x/oauth2"
	"wangzhiqiang/mpgrm/pkg/credential"
)

const (
	HOST        = "github.com"
	EnvPrefix   = "GITHUB"
	EnvUsername = EnvPrefix + credential.EnvUsernameSuffix
	EnvToken    = EnvPrefix + credential.EnvTokenSuffix
)

type Platform struct {
	Credential *credential.Credential
	ApiURL     string
}

func (p *Platform) GetClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.Credential.Token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func (p *Platform) WithCredential(credential *credential.Credential) error {
	if credential == nil || credential.Token == "" {
		return fmt.Errorf("invalid %s credential: credential is nil or token is empty", EnvPrefix)
	}
	p.Credential = credential
	return nil
}
