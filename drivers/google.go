package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/drivers/option"
	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/google"
)

const googleDriverName = "google"

func init() {
	err := RegisterDriver(
		option.Driver(googleDriverName),
		option.DefaultScopes(GoogleDefaultScopes),
		option.Callback(GoogleUserFn),
		option.Endpoint(google.Endpoint),
		option.APIMap(GoogleAPIMap),
		option.UserMap(GoogleUserMap),
	)
	if err != nil {
		panic(err)
	}
}

// GoogleUserMap is the map to create the User struct
var GoogleUserMap = map[string]string{
	"id":          "ID",
	"email":       "Email",
	"name":        "FullName",
	"given_name":  "FirstName",
	"family_name": "LastName",
	"picture":     "Avatar",
}

// GoogleAPIMap is the map for API endpoints
var GoogleAPIMap = map[string]string{
	"endpoint":     "https://www.googleapis.com",
	"userEndpoint": "/oauth2/v2/userinfo",
}

// GoogleUserFn is a callback to parse additional fields for User
var GoogleUserFn = func(client *http.Client, u *structs.User) {}

// GoogleDefaultScopes contains the default scopes
var GoogleDefaultScopes = []string{"profile", "email"}
