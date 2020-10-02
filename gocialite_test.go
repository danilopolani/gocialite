package gocialite

import (
    "testing"

    "github.com/gadelkareem/cachita"
    "github.com/stretchr/testify/assert"
)

var gocialTest = Gocial{c: cachita.Memory()}

func TestScopes(t *testing.T) {
    gocialTest.Scopes([]string{"email"})
    assert.Equal(t, gocialTest.ScopesArr, []string{"email"})
    assert.NotEqual(t, gocialTest.ScopesArr, []string{})

    gocialTest.
        Driver("google").
        Scopes([]string{"calendar.readonly"})
    assert.Equal(t, gocialTest.ScopesArr, []string{"profile", "email", "calendar.readonly"})
    assert.NotEqual(t, gocialTest.ScopesArr, []string{"profile", "email"})
    assert.NotEqual(t, gocialTest.ScopesArr, []string{})
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

    assert.Equal(gocialTest.Conf.ClientID, "foo")
    assert.NotEqual(gocialTest.Conf.ClientID, "")
    assert.NotNil(gocialTest.Conf.ClientID)

    assert.Equal(gocialTest.Conf.ClientSecret, "bar")
    assert.NotEqual(gocialTest.Conf.ClientSecret, "")
    assert.NotNil(gocialTest.Conf.ClientSecret)

    assert.Equal(gocialTest.Conf.RedirectURL, "http://example.com/auth/callback")
    assert.NotEqual(gocialTest.Conf.RedirectURL, "")
    assert.NotNil(gocialTest.Conf.RedirectURL)
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

    err = gocialTest.Handle(gocialTest.State, "foo")
    assert.NotNil(t, err)
}
