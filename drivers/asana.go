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
  userData := u.Raw["data"].(map[string]interface{})
  u.ID = fmt.Sprintf("%.0f", userData["id"].(float64))
  u.Email = userData["email"].(string)
  u.FullName = userData["name"].(string)

	// Set avatar
  if (userData["photo"] != nil) { 
	 u.Avatar = userData["photo"].(map[string]interface{})["image_1024x1024"].(string)
  }
}

// AsanaDefaultScopes contains the default scopes
var AsanaDefaultScopes = []string{}
