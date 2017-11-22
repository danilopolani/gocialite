package gocialite

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/linkedin"
	"gopkg.in/oleiade/reflections.v1"
)

// user struct
type user struct {
	ID        string                 `json:"id"`
	Username  string                 `json:"username"`
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	FullName  string                 `json:"full_name"`
	Email     string                 `json:"email"`
	Avatar    string                 `json:"avatar"`
	Raw       map[string]interface{} `json:"user"` // Raw data
}

// Gocial is the main struct of the package
type Gocial struct {
	driver, state string
	scopes        []string
	conf          *oauth2.Config
	User          user
}

// Set the basic information such as the endpoint and the scopes URIs
var mapAPI = map[string]map[string]string{
	"github": {
		"endpoint":       "https://api.github.com",
		"basicUserScope": "read:user",
		// Scopes
		"user":      "/user",
		"read:user": "/user",
		"repo":      "/user/repos",
	},
}

// Mapping to create a valid "user" struct from providers
var mapUser = map[string]map[string]string{
	"github": {
		"id":         "ID",
		"email":      "Email",
		"login":      "Username",
		"avatar_url": "Avatar",
		"name":       "FullName",
	},
}

// Driver is needed to choose the correct social
func (g *Gocial) Driver(driver string) *Gocial {
	g.driver = driver
	g.state = randToken()

	return g
}

// Scopes is used to set the oAuth scopes, for example "user", "calendar"
func (g *Gocial) Scopes(scopes []string) *Gocial {
	g.scopes = scopes
	return g
}

// Redirect returns an URL for the selected social oAuth login
func (g *Gocial) Redirect(clientID, clientSecret, redirectURL string) string {
	// Retrieve correct endpoint
	var endpoint oauth2.Endpoint
	switch g.driver {
	case "github":
		endpoint = github.Endpoint
	case "linkedin":
		endpoint = linkedin.Endpoint
	}

	g.conf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       g.scopes,
		Endpoint:     endpoint,
	}

	return g.conf.AuthCodeURL(g.state)
}

// Handle callback from provider
func (g *Gocial) Handle(state, code string) error {
	// Handle the exchange code to initiate a transport.
	if g.state != state {
		return fmt.Errorf("Invalid state: %s", state)
	}

	token, err := g.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return err
	}

	client := g.conf.Client(oauth2.NoContext, token)

	// Retrieve all from scopes
	driverAPI := mapAPI[g.driver]
	driverUserMap := mapUser[g.driver]
	scopes := g.scopes

	// Always append user scope
	if !inSlice("user", scopes) && !inSlice("read:user", scopes) {
		scopes = append([]string{driverAPI["basicUserScope"]}, scopes...)
	}

	for _, scope := range scopes {
		req, err := client.Get(driverAPI["endpoint"] + driverAPI[scope])
		if err != nil {
			return err
		}

		defer req.Body.Close()
		res, _ := ioutil.ReadAll(req.Body)

		// If the scope is about the user, save the details
		if scope == driverAPI["basicUserScope"] || scope == "user" {
			data, err := jsonDecode(res)
			if err != nil {
				return fmt.Errorf("Error decoding JSON: %s", err.Error())
			}

			// Scan all fields and dispatch through the mapping
			mapKeys := keys(driverUserMap)
			gUser := user{}
			for k, f := range data {
				if !inSlice(k, mapKeys) { // Skip if not in the mapping
					continue
				}

				// Assign the value
				_ = reflections.SetField(&gUser, driverUserMap[k], fmt.Sprint(f)) // Dirty way, but we need to convert also int/float to string
			}

			// Set the "raw" user interface
			gUser.Raw = data
			// Update the struct
			g.User = gUser
		}
	}

	return nil
}

// Generate a random token
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// Check if a value is in a string slice
func inSlice(v string, s []string) bool {
	for _, scope := range s {
		if scope == v {
			return true
		}
	}

	return false
}

// Decode a json or return an error
func jsonDecode(js []byte) (map[string]interface{}, error) {
	var decoded map[string]interface{}
	if err := json.Unmarshal(js, &decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}

func keys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
