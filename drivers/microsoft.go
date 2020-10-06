package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/microsoft"
)

const microsoftDriverName = "microsoft"

func init() {
	registerDriver(microsoftDriverName, MicrosoftDefaultScopes, MicrosoftUserFn, microsoft.AzureADEndpoint("common"), MicrosoftAPIMap, MicrosoftUserMap)
}

// MicrosoftUserMap is the map to create the User struct
var MicrosoftUserMap = map[string]string{
	"id":                "ID",
	"givenName":         "Username",
	"displayName":       "FullName",
	"userPrincipalName": "Email",
}

// MicrosoftAPIMap is the map for API endpoints
var MicrosoftAPIMap = map[string]string{
	"endpoint":     "https://graph.microsoft.com",
	"userEndpoint": "/v1.0/me",
}

// MicrosoftUserFn is a callback to parse additional fields for User
var MicrosoftUserFn = func(client *http.Client, u *structs.User) {}

// MicrosoftDefaultScopes contains the default scopes
var MicrosoftDefaultScopes = []string{"User.Read"}
