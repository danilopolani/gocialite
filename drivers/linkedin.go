package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/linkedin"
)

const (
	linkedinDriverName = "linkedin"
)

func init() {
	registerDriver(linkedinDriverName, LinkedInDefaultScopes, LinkedInUserFn, linkedin.Endpoint, LinkedInAPIMap, LinkedInUserMap)
}

// LinkedInUserMap is the map to create the User struct
var LinkedInUserMap = map[string]string{
	"id":            "ID",
	"vanityName":    "Username",
	"firstName":     "FirstName",
	"lastName":      "LastName",
	"formattedName": "FullName",
	"emailAddress":  "Email",
	"pictureUrl":    "Avatar",
}

// LinkedInAPIMap is the map for API endpoints
var LinkedInAPIMap = map[string]string{
	"endpoint":     "https://api.linkedin.com",
	"userEndpoint": "/v2/me",
	"emailEndpoint": "/v2/emailAddress?q=members&projection=(elements*(handle~))",
}

// LinkedInUserFn is a callback to parse additional fields for User
var LinkedInUserFn = func(client *http.Client, u *structs.User) {}

// LinkedInDefaultScopes contains the default scopes
var LinkedInDefaultScopes = []string{"r_emailaddress", "r_liteprofile", "w_member_social"}
