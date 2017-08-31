package controller

import (
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/ChrisTheBaron/strava-ical/model"
	"github.com/ChrisTheBaron/strava-ical/services"
	"github.com/ChrisTheBaron/strava-ical/utils/ical"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

type Calendar struct {
	controller
	cm *model.Calendar
	um *model.User
	sf services.StravaFactory
}

// NewCalendar returns a new controller.Calendars
func NewCalendar(
	con *entities.Config,
	cm *model.Calendar,
	um *model.User,
	sf services.StravaFactory) *Calendar {
	return &Calendar{controller{config: con}, cm, um, sf}
}

// Get lists the calendars that belong to the currently logged in user.
func (c *Calendar) Get(w http.ResponseWriter, r *http.Request) {

	uid, err := c.getUserIdFromContext(r)

	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry, something went wrong."))
		return
	}

	cals, err := c.cm.GetAllForUser(uid)

	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry, something went wrong."))
		return
	}

	fmt.Fprintf(w, "Authenticated as user id: %d\n\n", uid)
	fmt.Fprintln(w, "Users calendars:\n\n")

	for _, cal := range cals {
		fmt.Fprintf(w, "UUID: %s\n", cal.GetId().String())
	}

}

func (c *Calendar) GetICALById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	if sid, ok := vars["id"]; ok {

		id, err := uuid.FromString(sid)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Id"))
			return
		}

		uid, err := c.getUserIdFromContext(r)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Not logged in"))
			return
		}

		cal, err := c.cm.GetById(id)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			glog.Error(err)
			return
		}

		if cal.GetUserId() != uid {
			w.WriteHeader(http.StatusForbidden)
			glog.Warningln("Calendar not owned by current user")
			return
		}

		user, err := c.um.GetById(uid)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			glog.Error(err)
			return
		}

		strava := c.sf.NewStrava(user.GetStravaAccessToken(), user.GetStravaId())

		acts, err := strava.GetUsersActivities()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			glog.Error(err)
			return
		}

		ic := c.config.Calendar

		ic.NAME = fmt.Sprintf("Strava Activities - %s %s", user.GetFirstname(), user.GetLastname())
		ic.X_WR_CALNAME = ic.NAME

		for _, act := range acts {

			ic.AddComponent(ical.VComponent(ical.VEvent{
				UID:         fmt.Sprintf("%d", act.Id),
				SUMMARY:     act.Name,
				DESCRIPTION: fmt.Sprintf("%s\n\nDistance: %.2fkm", act.Name, act.Distance/1000),
				DTSTART:     act.StartDateLocal,
				DTEND:       act.StartDateLocal.Add(time.Duration(act.ElapsedTime) * time.Second),
				DTSTAMP:     act.StartDateLocal,
				LOCATION:    act.City,
				TZID:        "Europe/London",
				AllDay:      false,
			}))

		}

		w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		w.Header().Set("Content-Disposition", "inline; filename=strava.ics")
		ic.Encode(w)

	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No Id set"))
		return
	}

}

func (c *Calendar) GetById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	if sid, ok := vars["id"]; ok {

		id, err := uuid.FromString(sid)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Id"))
			return
		}

		uid, err := c.getUserIdFromContext(r)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Not logged in"))
			return
		}

		cal, err := c.cm.GetById(id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			glog.Error(err)
			return
		}

		if cal.GetUserId() != uid {
			w.WriteHeader(http.StatusForbidden)
			glog.Warningln("Calendar not owned by current user")
			return
		}

		fmt.Fprintf(w, "Calendar ID: %s", id.String())

	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No Id set"))
		return
	}

}

// Post inserts a new calendar in the db and redirects to it on success
func (c *Calendar) Post(w http.ResponseWriter, r *http.Request) {

	uid, err := c.getUserIdFromContext(r)

	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry, something went wrong."))
		return
	}

	id := uuid.NewV4()

	cal := entities.NewCalendar(id, uid)

	insertId, err := c.cm.Insert(cal)

	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry, something went wrong."))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", c.config.Slugs.Calendars, insertId.String()), http.StatusCreated)

}
