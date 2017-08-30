package utils

import (
	"github.com/BurntSushi/toml"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"io/ioutil"
	"path/filepath"
)

// GetConfigFromFile reads the website config from the given path.
//
// path is a filepath, relative to the current working directory, of a
// TOML file marshallable to a structs.Config struct.
//
// Returns a config struct and nil if the config read was successful,
// and an undefined value and non-nil otherwise.
func GetConfigFromFile(path string) (c entities.Config, err error) {
	absPath, _ := filepath.Abs(path)
	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		return
	}
	s := string(b)
	_, err = toml.Decode(s, &c)
	return
}

func ValidateConfig(config entities.Config) []string {

	var errs []string

	if config.StravaClientId == 0 {
		errs = append(errs, "StravaClientId required")
	}

	if config.StravaClientSecret == "" {
		errs = append(errs, "StravaClientSecret required")
	}

	return errs

}
