package gocialite

type Provider struct {
	Name         string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

func (p *Provider) SetScopes(scopes []string) *Provider {
	p.Scopes = scopes

	return p
}
