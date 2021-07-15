package gocialite

import (
	"sync"

	"golang.org/x/oauth2"
)

type Gocial struct {
	driver, state string
	scopes        []string
	conf          *oauth2.Config
	User          User
	Token         *oauth2.Token
}

type dispatcher struct {
	mu sync.RWMutex
	g  map[string]*Gocial
}

var providers = []Provider{}

func Create(providers ...*Provider) *Gocial {
	return &Gocial{}
}

func UseProvider(name string, clientID string, clientSecret string, redirectURL string) *Provider {
	return &Provider{
		name,
		clientID,
		clientSecret,
		redirectURL,
		[]string{},
	}
}
