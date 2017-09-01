package entities

import (
	"github.com/ChrisTheBaron/strava-ical/utils/ical"
)

// Config is a structure containing global website configuration.
type Config struct {
	Protocol           string         `toml:"protocol"` // http or https
	RootUrl            string         `toml:"rootUrl"`
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
	ListenAddress string `toml:"listenAddress"`
	ListenPort    int    `toml:"listenPort"`
	Timeout       int    `toml:"timeout"`
}

type Slugs struct {
	Calendars     string
	OAuth         string
	OAuthCallback string
	E404          string
	E500          string
}
