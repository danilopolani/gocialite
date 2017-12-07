package drivers

import (
	"encoding/json"
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
)

var (
	initAPIMap           = map[string]map[string]string{}
	initUserMap          = map[string]map[string]string{}
	initEndpointMap      = map[string]oauth2.Endpoint{}
	initCallbackMap      = map[string]func(client *http.Client, u *structs.User){}
	initDefaultScopesMap = map[string][]string{}
)

func registerDriver(driver string, defaultscopes []string, callback func(client *http.Client, u *structs.User), endpoint oauth2.Endpoint, apimap, usermap map[string]string) {
	initAPIMap[driver] = apimap
	initUserMap[driver] = usermap
	initEndpointMap[driver] = endpoint
	initCallbackMap[driver] = callback
	initDefaultScopesMap[driver] = defaultscopes
}

// InitializeDrivers adds all the drivers to the register func
func InitializeDrivers(register func(driver string, defaultscopes []string, callback func(client *http.Client, u *structs.User), endpoint oauth2.Endpoint, apimap, usermap map[string]string)) {
	for k := range initAPIMap {
		register(k, initDefaultScopesMap[k], initCallbackMap[k], initEndpointMap[k], initAPIMap[k], initUserMap[k])
	}
}

// Decode a json or return an error
func jsonDecode(js []byte) (map[string]interface{}, error) {
	var decoded map[string]interface{}
	if err := json.Unmarshal(js, &decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}
