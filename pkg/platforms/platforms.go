package platforms

import (
	"fmt"
	"github.com/chihqiang/mpgrm/pkg/credential"
	"net/url"
)

var platforms = map[string]func() IPlatform{}

func Register(host, envPrefix string, platform func() IPlatform) {
	platforms[host] = platform
	credential.Register(host, envPrefix)
}
func GetPlatform(repo *url.URL, cred *credential.Credential) (IPlatform, error) {
	newProviderFunc, ok := platforms[repo.Host]
	if !ok {
		return nil, fmt.Errorf("unsupported repository platform: %s", repo.Host)
	}
	provider := newProviderFunc()
	if err := provider.WithCredential(cred); err != nil {
		return nil, err
	}
	return provider, nil
}
