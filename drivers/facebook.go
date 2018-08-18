package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/drivers/option"
	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/facebook"
)

const facebookDriverName = "facebook"

func init() {
	err := RegisterDriver(
		option.Driver(facebookDriverName),
		option.DefaultScopes(FacebookDefaultScopes),
		option.Callback(FacebookUserFn),
		option.Endpoint(facebook.Endpoint),
		option.APIMap(FacebookAPIMap),
		option.UserMap(FacebookUserMap),
	)
	if err != nil {
		panic(err)
	}
}

// FacebookUserMap is the map to create the User struct
var FacebookUserMap = map[string]string{
	"id":         "ID",
	"email":      "Email",
	"name":       "FullName",
	"first_name": "FirstName",
	"last_name":  "LastName",
}

// FacebookAPIMap is the map for API endpoints
var FacebookAPIMap = map[string]string{
	"endpoint":     "https://graph.facebook.com",
	"userEndpoint": "/me?fields=id,name,first_name,last_name,email",
}

// FacebookUserFn is a callback to parse additional fields for User
var FacebookUserFn = func(client *http.Client, u *structs.User) {
	u.Avatar = FacebookAPIMap["endpoint"] + "/v2.8/" + u.ID + "/picture?width=800"
}

// FacebookDefaultScopes contains the default scopes
var FacebookDefaultScopes = []string{"email"}
