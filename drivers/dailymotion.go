package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
)

const dailymotionDriverName = "dailymotion"

func init() {
	registerDriver(dailymotionDriverName, DailyMotionDefaultScopes, DailyMotionUserFn, DailyMotionEndpoint, DailyMotionAPIMap, DailyMotionUserMap)
}

// DailyMotionEndpoint is the oAuth endpoint
var DailyMotionEndpoint = oauth2.Endpoint{
	AuthURL:  "https://www.dailymotion.com/oauth/authorize",
	TokenURL: "https://api.dailymotion.com/oauth/token",
}

// DailyMotionUserMap is the map to create the User struct
var DailyMotionUserMap = map[string]string{
	"id":             "ID",
	"username":       "Username",
	"fullname":       "FullName",
	"first_name":     "FirstName",
	"last_name":      "LastName",
	"email":          "Email",
	"avatar_720_url": "Avatar",
}

// DailyMotionAPIMap is the map for API endpoints
var DailyMotionAPIMap = map[string]string{
	"endpoint":     "https://api.dailymotion.com",
	"userEndpoint": "/me",
}

// DailyMotionUserFn is a callback to parse additional fields for User
var DailyMotionUserFn = func(client *http.Client, u *structs.User) {}

// DailyMotionDefaultScopes contains the default scopes
var DailyMotionDefaultScopes = []string{}
