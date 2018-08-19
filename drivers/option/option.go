// Package option implements Functional Options idiom
package option

import (
	"fmt"
	"net/http"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
)

// CallbackFunc is callback to parse additional fields for User
type CallbackFunc func(client *http.Client, u *structs.User)

// SigFunc is a function for signature calculation
type SigFunc func(
	conf *oauth2.Config,
	tok *oauth2.Token,
	params map[string]string,
) string

// RegisterFunc is function for driver registration
type RegisterFunc func(info Options)

// Setter of option
type Setter func(info *Options)

// Options of driver
type Options struct {
	driver        string
	defaultScopes []string
	callback      CallbackFunc
	endpoint      oauth2.Endpoint
	apiMap        map[string]string
	userMap       map[string]string
	sig           SigFunc
}

// Validate checks driver
func (o Options) Validate() error {
	if len(o.driver) == 0 {
		return fmt.Errorf("Driver not valid: empty driver name")
	}
	if len(o.apiMap) == 0 {
		return fmt.Errorf("Driver %s not valid: empty apiMap", o.driver)
	}
	// TODO: validate others
	return nil
}

/*

	Getters

*/

// Driver name
func (o Options) Driver() string { return o.driver }

// DefaultScopes of driver
func (o Options) DefaultScopes() []string { return o.defaultScopes }

// Callback to parse additional fields for User
func (o Options) Callback() CallbackFunc { return o.callback }

// Endpoint of driver
func (o Options) Endpoint() oauth2.Endpoint { return o.endpoint }

// APIMap of driver
func (o Options) APIMap() map[string]string { return o.apiMap }

// UserMap of driver
func (o Options) UserMap() map[string]string { return o.userMap }

// Sig of driver
func (o Options) Sig() SigFunc { return o.sig }

/*

	Setters

*/

// Driver option setter
func Driver(driver string) Setter {
	return func(info *Options) { info.driver = driver }
}

// DefaultScopes option setter
func DefaultScopes(defaultScopes []string) Setter {
	return func(info *Options) { info.defaultScopes = defaultScopes }
}

// Callback option setter
func Callback(callback CallbackFunc) Setter {
	return func(info *Options) { info.callback = callback }
}

// Endpoint option setter
func Endpoint(endpoint oauth2.Endpoint) Setter {
	return func(info *Options) { info.endpoint = endpoint }
}

// APIMap option setter
func APIMap(apiMap map[string]string) Setter {
	return func(info *Options) { info.apiMap = apiMap }
}

// UserMap option setter
func UserMap(userMap map[string]string) Setter {
	return func(info *Options) { info.userMap = userMap }
}

// Sig option setter
// you can use %SIG in a userEndpoint
func Sig(sig SigFunc) Setter {
	return func(info *Options) { info.sig = sig }
}
