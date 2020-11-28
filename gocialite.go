package gocialite

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/redaLaanait/gocialite/v2/drivers"
	"github.com/redaLaanait/gocialite/v2/stores"
	"github.com/redaLaanait/gocialite/v2/structs"
	"golang.org/x/oauth2"
	"gopkg.in/oleiade/reflections.v1"
)

// Dispatcher allows to safely issue concurrent Gocials
type Dispatcher struct {
	store stores.GocialStore
}

// NewDispatcher creates new Dispatcher
func NewDispatcher(store stores.GocialStore) *Dispatcher {
	if store == nil {
		store = stores.NewMemoryStore()
	}
	return &Dispatcher{store: store}
}

func encodeGocial(g *Gocial) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(g); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeGocial(b []byte, g *Gocial) error {
	buf := bytes.NewBuffer(b)
	enc := gob.NewDecoder(buf)

	return enc.Decode(g)
}

// New Gocial instance
func (d *Dispatcher) New() (*Gocial, error) {
	state := randToken()
	g := &Gocial{state: state}

	b, err := encodeGocial(g)
	if err != nil {
		return nil, fmt.Errorf("encode Gocial failed: %w", err)
	}
	if err := d.store.Save(state, b); err != nil {
		return nil, fmt.Errorf("save state failed: %w", err)
	}

	return g, nil
}

// Handle callback. Can be called only once for given state.
func (d *Dispatcher) Handle(state, code string) (*structs.User, *oauth2.Token, error) {
	b, err := d.store.Get(state)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid CSRF token: %s", state)
	}
	g := &Gocial{}
	if err = decodeGocial(b, g); err != nil {
		return nil, nil, err
	}

	err = g.Handle(state, code)
	err = d.store.Delete(state)

	return &g.User, g.Token, err
}

// Gocial is the main struct of the package
type Gocial struct {
	driver, state string
	scopes        []string
	conf          *oauth2.Config
	User          structs.User
	Token         *oauth2.Token
}

func init() {
	drivers.InitializeDrivers(RegisterNewDriver)
}

var (
	// Set the basic information such as the endpoint and the scopes URIs
	apiMap = map[string]map[string]string{}

	// Mapping to create a valid "user" struct from providers
	userMap = map[string]map[string]string{}

	// Map correct endpoints
	endpointMap = map[string]oauth2.Endpoint{}

	// Map custom callbacks
	callbackMap = map[string]func(client *http.Client, u *structs.User){}

	// Default scopes for each driver
	defaultScopesMap = map[string][]string{}
)

//RegisterNewDriver adds a new driver to the existing set
func RegisterNewDriver(driver string, defaultscopes []string, callback func(client *http.Client, u *structs.User), endpoint oauth2.Endpoint, apimap, usermap map[string]string) {
	apiMap[driver] = apimap
	userMap[driver] = usermap
	endpointMap[driver] = endpoint
	callbackMap[driver] = callback
	defaultScopesMap[driver] = defaultscopes
}

// Driver is needed to choose the correct social
func (g *Gocial) Driver(driver string) *Gocial {
	g.driver = driver
	g.scopes = defaultScopesMap[driver]

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

// Redirect returns an URL for the selected social oAuth login
func (g *Gocial) Redirect(clientID, clientSecret, redirectURL string) (string, error) {
	// Check if driver is valid
	if !inSlice(g.driver, complexKeys(apiMap)) {
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
		Endpoint:     endpointMap[g.driver],
	}

	return g.conf.AuthCodeURL(g.state), nil
}

// Handle callback from provider
func (g *Gocial) Handle(state, code string) error {
	// Handle the exchange code to initiate a transport.
	if g.state != state {
		return fmt.Errorf("Invalid state: %s", state)
	}

	// Check if driver is valid
	if !inSlice(g.driver, complexKeys(apiMap)) {
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
	driverAPIMap := apiMap[g.driver]
	driverUserMap := userMap[g.driver]
	userEndpoint := strings.Replace(driverAPIMap["userEndpoint"], "%ACCESS_TOKEN", token.AccessToken, -1)

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
	callbackMap[g.driver](client, &gUser)

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
	decoder := json.NewDecoder(strings.NewReader(string(js)))
	decoder.UseNumber()

	if err := decoder.Decode(&decoded); err != nil {
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
