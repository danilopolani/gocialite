package drivers

import (
	"net/http"
  "fmt"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
)

const asanaDriverName = "asana"

func init() {
	registerDriver(asanaDriverName, AsanaDefaultScopes, AsanaUserFn, AsanaEndpoint, AsanaAPIMap, AsanaUserMap)
}

// DailyMotionEndpoint is the oAuth endpoint
var AsanaEndpoint = oauth2.Endpoint{
  AuthURL:  "https://app.asana.com/-/oauth_authorize",
  TokenURL: "https://app.asana.com/-/oauth_token",
}

// AsanaUserMap is the map to create the User struct
var AsanaUserMap = map[string]string{}

// AsanaAPIMap is the map for API endpoints
var AsanaAPIMap = map[string]string{
	"endpoint":      "https://app.asana.com/api/1.0",
	"userEndpoint":  "/users/me?opt_fields=id,name,email,photo",
}

// AsanaUserFn is a callback to parse additional fields for User
var AsanaUserFn = func(client *http.Client, u *structs.User) {
  u.ID = fmt.Sprintf("%.0f", u.Raw["data"].(map[string]interface{})["id"].(float64))
  u.Email = u.Raw["data"].(map[string]interface{})["email"].(string)
  u.FullName = u.Raw["data"].(map[string]interface{})["name"].(string)

	// Set avatar
	u.Avatar = u.Raw["data"].(map[string]interface{})["photo"].(map[string]interface{})["image_1024x1024"].(string)
}

// AsanaDefaultScopes contains the default scopes
var AsanaDefaultScopes = []string{}
