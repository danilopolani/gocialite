package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/github"
)

const (
	githubdrivername = "github"
)

func init() {
	registerDriver(githubdrivername, GithubDefaultScopes, GithubUserFn, github.Endpoint, GithubAPIMap, GithubUserMap)
}

// GithubUserMap is the map to create the User struct
var GithubUserMap = map[string]string{
	"id":         "ID",
	"email":      "Email",
	"login":      "Username",
	"avatar_url": "Avatar",
	"name":       "FullName",
}

// GithubAPIMap is the map for API endpoints
var GithubAPIMap = map[string]string{
	"endpoint":     "https://api.github.com",
	"userEndpoint": "/user",
}

// GithubUserFn is a callback to parse additional fields for User
var GithubUserFn = func(client *http.Client, u *structs.User) {}

// GithubDefaultScopes contains the default scopes
var GithubDefaultScopes = []string{}
