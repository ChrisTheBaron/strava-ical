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

func (s *Strava) GetUsersActivities() (chan strava.ActivitySummary, chan error) {

	acts := make(chan strava.ActivitySummary, 200)
	err := make(chan error, 1)

	service := strava.NewAthletesService(s.client)

	call := service.ListActivities(s.stravaId)

	call.PerPage(200)

	var page int = 1

	go func() {
		for {
			actspage, e := call.Page(page).Do()
			if e != nil {
				err <- e
				close(err)
				close(acts)
				break
			}
			if len(actspage) == 0 {
				close(err)
				close(acts)
				break
			}
			for _, act := range actspage {
				acts <- *act
			}
			page++
		}
	}()

	return acts, err

}
