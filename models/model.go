package models

import (
	"github.com/ChrisTheBaron/strava-ical/structs"
	"github.com/strava/go.strava"
)

// ModelInterface is the interface to which models adhere.
type ModelInterface interface {
	Get() (data *interface{}, err error) // @TODO: Refactor this to a more appropriate name
}

// Model is the base type of models in the 2016site architecture.
//
// Each model wraps the MyRadio API to provide the data required for a
// controller.
type Model struct {
	client *strava.Client
	config *structs.Config
}
