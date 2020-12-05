// Package drivers is a singleton for drivers
package drivers

import (
	"encoding/json"

	"github.com/danilopolani/gocialite/drivers/option"
)

var drivers = make(map[string]*option.Options)

// RegisterDriver adds a new driver to the existing set
func RegisterDriver(setters ...option.Setter) error {
	info := option.Options{}
	for _, s := range setters {
		s(&info)
	}
	if err := info.Validate(); err != nil {
		return err
	}

	drivers[info.Driver()] = &info
	return nil
}

// Driver returns options of driver
func Driver(name string) (o option.Options, ok bool) {
	if opt, ok := drivers[name]; opt != nil && ok {
		return *opt, true
	}
	return option.Options{}, false
}

// MustDriver returns options of driver
func MustDriver(name string) option.Options {
	opt, _ := Driver(name)
	return opt
}

// Decode a json or return an error
func jsonDecode(js []byte) (map[string]interface{}, error) {
	var decoded map[string]interface{}
	if err := json.Unmarshal(js, &decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}
