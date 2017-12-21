package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/foursquare"
)

const foursquareDriverName = "foursquare"

func init() {
	registerDriver(foursquareDriverName, FoursquareDefaultScopes, FoursquareUserFn, foursquare.Endpoint, FoursquareAPIMap, FoursquareUserMap)
}

// FoursquareUserMap is the map to create the User struct
var FoursquareUserMap = map[string]string{}

// FoursquareAPIMap is the map for API endpoints
var FoursquareAPIMap = map[string]string{
	"endpoint":     "https://api.foursquare.com",
	"userEndpoint": "/v2/users/self?oauth_token=%ACCESS_TOKEN&v=20171220",
}

// FoursquareUserFn is a callback to parse additional fields for User
var FoursquareUserFn = func(client *http.Client, u *structs.User) {
	user := u.Raw["response"].(map[string]interface{})["user"].(map[string]interface{})

	u.ID = user["id"].(string)
	u.FirstName = user["firstName"].(string)
	u.LastName = user["lastName"].(string)
	u.FullName = u.FirstName + " " + u.LastName

	if email, ok := user["contact"].(map[string]interface{})["email"]; ok {
		u.Email = email.(string)
	}
	if avatarPrefix, ok := user["photo"].(map[string]interface{})["prefix"]; ok {
		if avatarSuffix, ok2 := user["photo"].(map[string]interface{})["suffix"]; ok2 {
			u.Avatar = avatarPrefix.(string) + "original" + avatarSuffix.(string)
		}
	}
}

// FoursquareDefaultScopes contains the default scopes
var FoursquareDefaultScopes = []string{}
