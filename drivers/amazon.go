package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/drivers/option"
	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/amazon"
)

const amazonDriverName = "amazon"

func init() {
	err := RegisterDriver(
		option.Driver(amazonDriverName),
		option.DefaultScopes(AmazonDefaultScopes),
		option.Callback(AmazonUserFn),
		option.Endpoint(amazon.Endpoint),
		option.APIMap(AmazonAPIMap),
		option.UserMap(AmazonUserMap),
	)
	if err != nil {
		panic(err)
	}
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
