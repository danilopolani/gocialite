package gocialite

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/danilopolani/gocialite/drivers"
	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/linkedin"
	"gopkg.in/oleiade/reflections.v1"
)

// Gocial is the main struct of the package
type Gocial struct {
	driver, state string
	scopes        []string
	conf          *oauth2.Config
	User          structs.User
	Token         *oauth2.Token
}

// Set the basic information such as the endpoint and the scopes URIs
var apiMap = map[string]map[string]string{
	"github":   drivers.GithubAPIMap,
	"linkedin": drivers.LinkedInAPIMap,
}

// Mapping to create a valid "user" struct from providers
var userMap = map[string]map[string]string{
	"github":   drivers.GithubUserMap,
	"linkedin": drivers.LinkedInUserMap,
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
func (g *Gocial) Redirect(clientID, clientSecret, redirectURL string) (string, error) {
	// Retrieve correct endpoint
	var endpoint oauth2.Endpoint
	switch g.driver {
	case "github":
		endpoint = github.Endpoint
	case "linkedin":
		endpoint = linkedin.Endpoint
	default:
		return "", fmt.Errorf("Driver not valid: %s", g.driver)
	}

	g.conf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       g.scopes,
		Endpoint:     endpoint,
	}

	return g.conf.AuthCodeURL(g.state), nil
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

	// Set gocial token
	g.Token = token

	// Retrieve all from scopes
	driverAPIMap := apiMap[g.driver]
	driverUserMap := userMap[g.driver]

	// Get user info
	req, err := client.Get(driverAPIMap["endpoint"] + driverAPIMap["userEndpoint"])
	if err != nil {
		return err
	}

	defer req.Body.Close()
	res, _ := ioutil.ReadAll(req.Body)
	data, err := jsonDecode(res)
	if err != nil {
		return fmt.Errorf("Error decoding JSON: %s", err.Error())
	}

	// Scan all fields and dispatch through the mapping
	mapKeys := keys(driverUserMap)
	gUser := structs.User{}
	for k, f := range data {
		if !inSlice(k, mapKeys) { // Skip if not in the mapping
			continue
		}

		// Assign the value
		// Dirty way, but we need to convert also int/float to string
		_ = reflections.SetField(&gUser, driverUserMap[k], fmt.Sprint(f))
	}

	// Set the "raw" user interface
	gUser.Raw = data
	// Update the struct
	g.User = gUser

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
