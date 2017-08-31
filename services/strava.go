package services

import (
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/strava/go.strava"
)

type Strava struct {
	clientSecret string
	clientId     int
	stravaId     int64
	accessToken  string
	client       *strava.Client
}

type StravaFactory struct {
	NewStrava func(accessToken string, athleteId int64) Strava
}

func NewStravaFactory(config *entities.Config) StravaFactory {
	return StravaFactory{func(accessToken string, athleteId int64) Strava {
		return Strava{
			config.StravaClientSecret,
			config.StravaClientId,
			athleteId,
			accessToken,
			strava.NewClient(accessToken),
		}
	}}
}

func (s *Strava) GetUsersActivities() (acts []*strava.ActivitySummary, err error) {

	service := strava.NewAthletesService(s.client)

	call := service.ListActivities(s.stravaId)

	call.PerPage(200)

	var page int = 1

	for {
		actspage, err := call.Page(page).Do()
		if err != nil || len(actspage) == 0 {
			break
		}
		acts = append(acts, actspage...)
		page++
	}

	return

}
