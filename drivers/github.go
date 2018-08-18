package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/drivers/option"
	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/github"
)

const githubDriverName = "github"

func init() {
	err := RegisterDriver(
		option.Driver(githubDriverName),
		option.DefaultScopes(GithubDefaultScopes),
		option.Callback(GithubUserFn),
		option.Endpoint(github.Endpoint),
		option.APIMap(GithubAPIMap),
		option.UserMap(GithubUserMap),
	)
	if err != nil {
		panic(err)
	}
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
