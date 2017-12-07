package drivers

import (
	"io/ioutil"
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/slack"
)

const slackDriverName = "slack"

func init() {
	registerDriver(slackDriverName, SlackDefaultScopes, SlackUserFn, slack.Endpoint, SlackAPIMap, SlackUserMap)
}

// SlackUserMap is the map to create the User struct
var SlackUserMap = map[string]string{
	"real_name":      "FullName",
	"first_name":     "FirstName",
	"last_name":      "LastName",
	"email":          "Email",
	"image_original": "Avatar",
}

// SlackAPIMap is the map for API endpoints
var SlackAPIMap = map[string]string{
	"endpoint":     "https://slack.com/api",
	"userEndpoint": "/users.profile.get",
	"authEndpoint": "/auth.test",
}

// SlackUserFn is a callback to parse additional fields for User
var SlackUserFn = func(client *http.Client, u *structs.User) {
	// Get user ID
	req, err := client.Get(SlackAPIMap["endpoint"] + SlackAPIMap["authEndpoint"])
	if err != nil {
		return
	}

	defer req.Body.Close()
	res, _ := ioutil.ReadAll(req.Body)
	data, err := jsonDecode(res)
	if err != nil {
		return
	}

	u.ID = data["user_id"].(string)

	// Fetch other user information
	userInfo := u.Raw["profile"].(map[string]interface{})
	u.Username = userInfo["display_name"].(string)
	u.FullName = userInfo["real_name"].(string)
	u.FirstName = userInfo["first_name"].(string)
	u.LastName = userInfo["last_name"].(string)
	u.Email = userInfo["email"].(string)
	u.Avatar = userInfo["image_original"].(string)
}

// SlackDefaultScopes contains the default scopes
var SlackDefaultScopes = []string{"users.profile:read"}
