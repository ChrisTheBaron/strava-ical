package models

import (
	"github.com/ChrisTheBaron/strava-ical/structs"
	"github.com/strava/go.strava"
)

// IndexModel is the model for the Index controller.
type IndexModel struct {
	Model
}

// NewIndexModel returns a new IndexModel
func NewIndexModel(cli *strava.Client, con *structs.Config) *IndexModel {
	return &IndexModel{Model{client: cli, config: con}}
}

// Get gets the data required for the Index controller from Strava.
//
// Otherwise, it returns undefined data and the error causing failure.
func (m *IndexModel) Get() ([]*strava.ActivitySummary, error) {

	service := strava.NewAthletesService(m.client)

	call := service.ListActivities(m.config.StravaAthleteId)

	return call.Do()

}
