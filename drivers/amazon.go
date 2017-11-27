package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/amazon"
)

const amazonDriverName = "amazon"

func init() {
	registerDriver(amazonDriverName, AmazonDefaultScopes, AmazonUserFn, amazon.Endpoint, AmazonAPIMap, AmazonUserMap)
}

// AmazonUserMap is the map to create the User struct
var AmazonUserMap = map[string]string{
	"user_id": "ID",
	"name":    "FullName",
	"email":   "Email",
}

// AmazonAPIMap is the map for API endpoints
var AmazonAPIMap = map[string]string{
	"endpoint":     "https://api.amazon.com",
	"userEndpoint": "/user/profile",
}

// AmazonUserFn is a callback to parse additional fields for User
var AmazonUserFn = func(client *http.Client, u *structs.User) {}

// AmazonDefaultScopes contains the default scopes
var AmazonDefaultScopes = []string{"profile"}
