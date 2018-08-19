package gocialite

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/danilopolani/gocialite/drivers"
	"github.com/danilopolani/gocialite/drivers/option"
	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
	"gopkg.in/oleiade/reflections.v1"
)

// Dispatcher allows to safely issue concurrent Gocials
type Dispatcher struct {
	mu sync.RWMutex
	g  map[string]*Gocial
}

// NewDispatcher creates new Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{g: make(map[string]*Gocial)}
}

// New Gocial instance
func (d *Dispatcher) New() *Gocial {
	d.mu.Lock()
	defer d.mu.Unlock()
	state := randToken()
	g := &Gocial{state: state}
	g.params = make(map[string]string)
	d.g[state] = g
	return g
}

// Handle callback. Can be called only once for given state.
func (d *Dispatcher) Handle(state, code string) (*structs.User, *oauth2.Token, error) {
	d.mu.RLock()
	g, ok := d.g[state]
	d.mu.RUnlock()
	if !ok {
		return nil, nil, fmt.Errorf("invalid CSRF token: %s", state)
	}
	err := g.Handle(state, code)
	d.mu.Lock()
	delete(d.g, state)
	d.mu.Unlock()
	return &g.User, g.Token, err
}

// Gocial is the main struct of the package
type Gocial struct {
	driver, state string
	scopes        []string
	conf          *oauth2.Config
	User          structs.User
	Token         *oauth2.Token
	params        map[string]string
}

// RegisterNewDriver adds a new driver to the existing set
// The function is deprecated, use drivers.RegisterDriver instead
func RegisterNewDriver(
	driver string,
	defaultscopes []string,
	callback func(client *http.Client, u *structs.User),
	endpoint oauth2.Endpoint,
	apimap map[string]string,
	usermap map[string]string,
) {
	drivers.RegisterDriver(
		option.APIMap(apimap),
		option.UserMap(usermap),
		option.Endpoint(endpoint),
		option.Callback(callback),
		option.DefaultScopes(defaultscopes),
	)
}

// Driver is needed to choose the correct social
func (g *Gocial) Driver(driver string) *Gocial {
	g.driver = driver
	g.scopes = drivers.MustDriver(driver).DefaultScopes()

	// BUG: sequential usage of single Gocial instance will have same CSRF token. This is serious security issue.
	// NOTE: Dispatcher eliminates this bug.
	if g.state == "" {
		g.state = randToken()
	}

	return g
}

// Scopes is used to set the oAuth scopes, for example "user", "calendar"
func (g *Gocial) Scopes(scopes []string) *Gocial {
	g.scopes = append(g.scopes, scopes...)
	return g
}

// Params is used to set additional parameters for driver
// for exapmle: APPLICATION_KEY for OK driver
func (g *Gocial) Params(params map[string]string) *Gocial {
	g.params = make(map[string]string)
	for p, v := range params {
		g.params["%"+p] = strings.ToUpper(v)
	}
	return g
}

// Redirect returns an URL for the selected social oAuth login
func (g *Gocial) Redirect(clientID, clientSecret, redirectURL string) (string, error) {
	drv, ok := drivers.Driver(g.driver)
	// Check if driver is valid
	if !ok {
		return "", fmt.Errorf("Driver not valid: %s", g.driver)
	}

	// Check if valid redirectURL
	_, err := url.ParseRequestURI(redirectURL)
	if err != nil {
		return "", fmt.Errorf("Redirect URL <%s> not valid: %s", redirectURL, err.Error())
	}
	if !strings.HasPrefix(redirectURL, "http://") && !strings.HasPrefix(redirectURL, "https://") {
		return "", fmt.Errorf("Redirect URL <%s> not valid: protocol not valid", redirectURL)
	}

	g.conf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       g.scopes,
		Endpoint:     drv.Endpoint(),
	}

	return g.conf.AuthCodeURL(g.state), nil
}

// Handle callback from provider
func (g *Gocial) Handle(state, code string) error {
	drv, ok := drivers.Driver(g.driver)

	// Handle the exchange code to initiate a transport.
	if g.state != state {
		return fmt.Errorf("Invalid state: %s", state)
	}

	// Check if driver is valid
	if !ok {
		return fmt.Errorf("Driver not valid: %s", g.driver)
	}

	token, err := g.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return fmt.Errorf("oAuth exchanged failed: %s", err.Error())
	}

	client := g.conf.Client(oauth2.NoContext, token)

	// Set gocial token
	g.Token = token

	// Retrieve all from scopes
	driverAPIMap := drv.APIMap()
	driverUserMap := drv.UserMap()
	userEndpoint := driverAPIMap["userEndpoint"]
	g.params["%ACCESS_TOKEN"] = token.AccessToken
	if drv.Sig() != nil {
		g.params["%SIG"] = drv.Sig()(g.conf, g.Token, g.params)
	}
	for p, v := range g.params {
		userEndpoint = strings.Replace(userEndpoint, p, v, -1)
	}

	// Get user info
	req, err := client.Get(driverAPIMap["endpoint"] + userEndpoint)
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

	// Custom callback
	drv.Callback()(client, &gUser)

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

// Return the keys of a map
func keys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func complexKeys(m map[string]map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
