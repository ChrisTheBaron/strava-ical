package controller

import (
	"errors"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/ChrisTheBaron/strava-ical/middleware"
	"net/http"
)

type controller struct {
	config *entities.Config
}

func (c *controller) getUserIdFromContext(r *http.Request) (int, error) {

	uid, ok := r.Context().Value(middleware.UserIdCtxKey).(int)

	if !ok {
		return 0, errors.New("No user ID set in context")
	}

	if uid == 0 {
		return 0, errors.New("No user ID set in context")
	}

	return uid, nil

}
