package drivers

import (
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
