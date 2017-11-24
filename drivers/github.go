package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/structs"
)

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
