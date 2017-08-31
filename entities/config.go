package entities

import (
	"github.com/ChrisTheBaron/strava-ical/utils/ical"
)

// Config is a structure containing global website configuration.
type Config struct {
	DBPath             string         `toml:"dbPath"`
	Slugs              Slugs          `toml:"slugs"`
	JWTCookieName      string         `toml:"jwtCookieName"`
	JWTKey             string         `toml:"jwtKey"`
	StravaClientId     int            `toml:"stravaClientId"`
	StravaClientSecret string         `toml:"stravaClientSecret"`
	Server             Server         `toml:"server"`
	Calendar           ical.VCalendar `toml:"calendar"`
}

// Server is a structure containing server configuration.
type Server struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
	Timeout int    `toml:"timout"`
}

type Slugs struct {
	Calendars     string
	Login         string
	OAuth         string
	OAuthCallback string
}
