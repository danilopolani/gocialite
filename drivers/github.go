package drivers

import (
	"encoding/json"
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2/github"
)

const githubDriverName = "github"

func init() {
	registerDriver(githubDriverName, GithubDefaultScopes, GithubUserFn, github.Endpoint, GithubAPIMap, GithubUserMap)
}

// GithubUserMap is the map to create the User struct
var GithubUserMap = map[string]string{
	"id":         "ID",
	"email":      "Email",
	"login":      "Username",
	"avatar_url": "Avatar",
	"name":       "FullName",
}

// GithubAPIMap is the map for API endpoints
var GithubAPIMap = map[string]string{
	"endpoint":      "https://api.github.com",
	"userEndpoint":  "/user",
	"emailEndpoint": "/user/emails",
}

// GithubUserFn is a callback to parse additional fields for User
var GithubUserFn = func(client *http.Client, u *structs.User) {
	// Used to parse the email from response
	type additionalEmail struct {
		Email string `json:"email"`
	}
	var email []additionalEmail

	// Email can be nil because of the "keep my email private" setting
	if u.Email == "<nil>" {
		// Retrieve email
		req, err := client.Get(GithubAPIMap["endpoint"] + GithubAPIMap["emailEndpoint"])
		if err != nil {
			return
		}

		defer req.Body.Close()
		err = json.NewDecoder(req.Body).Decode(&email)
		if err != nil {
			return
		}

		u.Email = email[0].Email
	}
}

// GithubDefaultScopes contains the default scopes
var GithubDefaultScopes = []string{"user:email"}
