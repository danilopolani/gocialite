package drivers

import (
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
)

var (
	initapiMap           = map[string]map[string]string{}
	inituserMap          = map[string]map[string]string{}
	initendpointMap      = map[string]oauth2.Endpoint{}
	initcallbackMap      = map[string]func(client *http.Client, u *structs.User){}
	initdefaultScopesMap = map[string][]string{}
)

// func init() {
// 	initapiMap = make(map[string]map[string]string)
// 	inituserMap = make(map[string]map[string]string)
// 	initendpointMap = make(map[string]oauth2.Endpoint)
// 	initcallbackMap = make(map[string]func(client *http.Client, u *structs.User))
// 	initdefaultScopesMap = make(map[string][]string)
// }

func registerDriver(driver string, defaultscopes []string, callback func(client *http.Client, u *structs.User), endpoint oauth2.Endpoint, apimap, usermap map[string]string) {
	initapiMap[driver] = apimap
	inituserMap[driver] = usermap
	initendpointMap[driver] = endpoint
	initcallbackMap[driver] = callback
	initdefaultScopesMap[driver] = defaultscopes
}

//InitializeDrivers adds all the drivers to the register func
func InitializeDrivers(register func(driver string, defaultscopes []string, callback func(client *http.Client, u *structs.User), endpoint oauth2.Endpoint, apimap, usermap map[string]string)) {
	for k := range initapiMap {
		register(k, initdefaultScopesMap[k], initcallbackMap[k], initendpointMap[k], initapiMap[k], inituserMap[k])
	}
}
