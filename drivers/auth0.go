package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
)

const auth0DriverName = "auth0"

func init() {
	registerDriver(auth0DriverName, Auth0DefaultScopes, Auth0UserFn, Auth0Endpoint, Auth0APIMap, Auth0UserMap)
}

// Auth0Endpoint is the oAuth endpoint
var Auth0Endpoint = oauth2.Endpoint{
	AuthURL:  "/authorize",
	TokenURL: "/oauth/token",
}

// Auth0UserMap is the map to create the User struct
var Auth0UserMap = map[string]string{
	"user_id":  "ID",
	"nickname": "Username",
	"name":     "FullName",
	"email":    "Email",
	"picture":  "Avatar",
}

// Auth0APIMap is the map for API endpoints
var Auth0APIMap = map[string]string{
	"endpoint":     "https://api.Auth0.com",
	"userEndpoint": "/userinfo",
}

// Auth0UserFn is a callback to parse additional fields for User
var Auth0UserFn = func(client *http.Client, u *structs.User) {}

// Auth0DefaultScopes contains the default scopes
var Auth0DefaultScopes = []string{}
