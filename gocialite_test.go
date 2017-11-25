package gocialite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var gocialTest Gocial

func TestScopes(t *testing.T) {
	gocialTest.Scopes([]string{"email"})
	assert.Equal(t, gocialTest.scopes, []string{"email"})
	assert.NotEqual(t, gocialTest.scopes, []string{})
}
func TestConf(t *testing.T) {
	assert := assert.New(t)

	gocialTest.
		Driver("github").
		Redirect(
			"foo",
			"bar",
			"http://example.com/auth/callback",
		)

	assert.Equal(gocialTest.conf.ClientID, "foo")
	assert.NotEqual(gocialTest.conf.ClientID, "")
	assert.NotNil(gocialTest.conf.ClientID)

	assert.Equal(gocialTest.conf.ClientSecret, "bar")
	assert.NotEqual(gocialTest.conf.ClientSecret, "")
	assert.NotNil(gocialTest.conf.ClientSecret)

	assert.Equal(gocialTest.conf.RedirectURL, "http://example.com/auth/callback")
	assert.NotEqual(gocialTest.conf.RedirectURL, "")
	assert.NotNil(gocialTest.conf.RedirectURL)
}
func TestDriver(t *testing.T) {
	var err error

	_, err = gocialTest.Driver("unknown").
		Redirect(
			"xxxxxxxx",
			"xxxxxxxxxxxxxxxxxxxxxxxx",
			"http://example.com/auth/callback",
		)
	assert.NotNil(t, err)

	_, err = gocialTest.Driver("github").
		Redirect(
			"xxxxxxxx",
			"xxxxxxxxxxxxxxxxxxxxxxxx",
			"http://example.com/auth/callback",
		)
	assert.Nil(t, err)
}
func TestRedirectURL(t *testing.T) {
	var err error

	_, err = gocialTest.Driver("github").
		Redirect(
			"xxxxxxxx",
			"xxxxxxxxxxxxxxxxxxxxxxxx",
			"/auth/callback",
		)
	assert.NotNil(t, err)

	_, err = gocialTest.Driver("github").
		Redirect(
			"xxxxxxxx",
			"xxxxxxxxxxxxxxxxxxxxxxxx",
			"http://example.com/auth/callback",
		)
	assert.Nil(t, err)
}

func TestState(t *testing.T) {
	var err error

	err = gocialTest.Driver("github").
		Handle("fakeState", "foo")
	assert.NotNil(t, err)
}

func TestExchange(t *testing.T) {
	var err error

	// Generate a state
	gocialTest.
		Driver("github").
		Redirect(
			"xxxxxxxx",
			"xxxxxxxxxxxxxxxxxxxxxxxx",
			"http://example.com/auth/callback",
		)

	err = gocialTest.Handle(gocialTest.state, "foo")
	assert.NotNil(t, err)
}
