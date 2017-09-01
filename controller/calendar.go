package controller

import (
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/ChrisTheBaron/strava-ical/model"
	"github.com/ChrisTheBaron/strava-ical/services"
	"github.com/ChrisTheBaron/strava-ical/utils"
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
	tr *utils.Template
}

// NewCalendar returns a new controller.Calendars
func NewCalendar(
	con *entities.Config,
	cm *model.Calendar,
	um *model.User,
	sf services.StravaFactory,
	tr *utils.Template) *Calendar {
	return &Calendar{controller{config: con}, cm, um, sf, tr}
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

	user, err := c.um.GetById(uid)

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

	data := struct {
		User      entities.User
		Calendars []entities.Calendar
	}{
		user,
		cals,
	}

	err = c.tr.Render(w, data, "calendars.tmpl")

	if err != nil {
		glog.Error(err)
	}

}

func (c *Calendar) GetICALById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	sid, ok := vars["id"]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		glog.Error("No Id set")
		return
	}

	id, err := uuid.FromString(sid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		glog.Error(err)
		return
	}

	cal, err := c.cm.GetById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		glog.Error(err)
		return
	}

	user, err := c.um.GetById(cal.UserId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		glog.Error(err)
		return
	}

	strava := c.sf.NewStrava(user.StravaAccessToken, user.StravaId)

	acts, err := strava.GetUsersActivities()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		glog.Error(err)
		return
	}

	ic := c.config.Calendar

	ic.NAME = fmt.Sprintf("Strava Activities - %s %s", user.Firstname, user.Lastname)
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
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s.ics", id.String()))
	ic.Encode(w)

}

func (c *Calendar) GetById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	sid, ok := vars["id"]

	if !ok {
		glog.Warningln("No Id set")
		c.redirect(w, c.config.Slugs.Calendars)
		return
	}

	id, err := uuid.FromString(sid)

	if err != nil {
		glog.Warningln("Bad Id set")
		c.redirect(w, c.config.Slugs.E404)
		return
	}

	uid, err := c.getUserIdFromContext(r)

	if err != nil {
		glog.Warningln("Not logged in")
		c.redirect(w, "/")
		return
	}

	cal, err := c.cm.GetById(id)

	if err != nil {
		glog.Warningln("Bad Id set")
		c.redirect(w, c.config.Slugs.E500)
		return
	}

	if cal.UserId != uid {
		glog.Warningln("Calendar not owned by current user")
		c.redirect(w, c.config.Slugs.E404)
		return
	}

	data := struct {
		Calendar entities.Calendar
	}{
		Calendar: cal,
	}

	err = c.tr.Render(w, data, "calendar.tmpl")

	if err != nil {
		glog.Error(err)
	}

}

// Post inserts a new calendar in the db and redirects to it on success
func (c *Calendar) Post(w http.ResponseWriter, r *http.Request) {

	uid, err := c.getUserIdFromContext(r)

	if err != nil {
		glog.Error(err)
		c.redirect(w, c.config.Slugs.E500)
		return
	}

	id := uuid.NewV4()

	cal := entities.Calendar{Id: id, UserId: uid}

	insertId, err := c.cm.Insert(cal)

	if err != nil {
		glog.Error(err)
		c.redirect(w, c.config.Slugs.E500)
		return
	}

	c.redirectWithStatus(w, fmt.Sprintf("%s/%s", c.config.Slugs.Calendars, insertId.String()), http.StatusCreated)

}
