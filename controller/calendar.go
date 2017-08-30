package controller

import (
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/golang/glog"
	"net/http"
)

type Calendar struct {
	controller
}

// NewCalendar returns a new controller.Calendar
func NewCalendar(con *entities.Config) *Calendar {
	return &Calendar{controller{config: con}}
}

// Get lists the calendars that belong to the currently logged in user.
func (a *Calendar) Get(w http.ResponseWriter, r *http.Request) {

	uid, err := a.getUserIdFromContext(r)

	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry, something went wrong."))
		return
	}

	fmt.Fprintf(w, "Authenticated as user id: %d", uid)

}
