package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/drivers/option"
	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/linkedin"
)

const (
	linkedinDriverName = "linkedin"
)

func init() {
	err := RegisterDriver(
		option.Driver(linkedinDriverName),
		option.DefaultScopes(LinkedInDefaultScopes),
		option.Callback(LinkedInUserFn),
		option.Endpoint(linkedin.Endpoint),
		option.APIMap(LinkedInAPIMap),
		option.UserMap(LinkedInUserMap),
	)
	if err != nil {
		panic(err)
	}
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
	"userEndpoint": "/v1/people/~:(id,first-name,last-name,formatted-name,email-address,picture-url,maiden-name,headline,location,industry,current-share,num-connections,summary,specialties,positions,public-profile-url)?format=json",
}

// LinkedInUserFn is a callback to parse additional fields for User
var LinkedInUserFn = func(client *http.Client, u *structs.User) {}

// LinkedInDefaultScopes contains the default scopes
var LinkedInDefaultScopes = []string{}
