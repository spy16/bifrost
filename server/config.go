package server

import (
	"math/rand"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Config represents the server configuration.
type Config struct {
	Providers map[string]ProviderConfig `toml:"providers"`
}

// ProviderConfig represents the OAuth2 provider configuration.
type ProviderConfig struct {
	AuthURI        string            `toml:"auth_uri"`
	TokenURI       string            `toml:"token_uri"`
	ClientID       string            `toml:"client_id"`
	ClientSecret   string            `toml:"client_secret"`
	RedirectURL    string            `toml:"redirect_uri"`
	Scopes         []string          `toml:"scopes"`
	EndpointParams map[string]string `toml:"params"`
}

// AuthorizeURL generates the OAuth2 authorization URL for redirection.
func (pc *ProviderConfig) AuthorizeURL(state string) string {
	params := pc.ExtraParams()
	options := []oauth2.AuthCodeOption{}
	for k := range params {
		options = append(options, oauth2.SetAuthURLParam(k, params.Get(k)))
	}

	return pc.OAuth2().AuthCodeURL(state, options...)
}

// OAuth2 returns the OAuth2 client configuration object for the provider.
func (pc *ProviderConfig) OAuth2() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     pc.ClientID,
		ClientSecret: pc.ClientSecret,
		Scopes:       pc.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:   pc.AuthURI,
			TokenURL:  pc.TokenURI,
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		RedirectURL: pc.RedirectURL,
	}
}

// ClientCreds returns the client-credentials grant config object for the
// provider.
func (pc *ProviderConfig) ClientCreds() *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:       pc.ClientID,
		ClientSecret:   pc.ClientSecret,
		TokenURL:       pc.TokenURI,
		EndpointParams: pc.ExtraParams(),
		AuthStyle:      oauth2.AuthStyleInHeader,
		Scopes:         pc.Scopes,
	}
}

// ExtraParams returns any custom parameters defined for the provider.
func (pc *ProviderConfig) ExtraParams() url.Values {
	params := url.Values{}
	for k, v := range pc.EndpointParams {
		params.Set(k, v)
	}

	return params
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
