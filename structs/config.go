package structs

import "github.com/ChrisTheBaron/strava-ical/utils/ical"

// Config is a structure containing global website configuration.
//
// See the comments for Server and PageContext for more details.
type Config struct {
	StravaApiKey             string         `toml:"stravaApiKey"`
	StravaAthleteId          int64          `toml:"stravaAthleteId"`
	Server                   Server         `toml:"server"`
	Calendar                 ical.VCalendar `toml:"calendar"`
	CalendarEntryName        string         `toml:"calendarEntryName"`
	CalendarEntryDescription string         `toml:"calendarEntryDescription"`
	CalendarDescription      string         `toml:"calendarDescription"`
	ShortName                string         `toml:"shortName"`
	LongName                 string         `toml:"longName"`
}

// Server is a structure containing server configuration.
type Server struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
	Timeout int    `toml:"timout"`
}
