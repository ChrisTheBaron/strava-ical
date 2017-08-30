package controllers

import (
	"bytes"
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/structs"
	"github.com/ChrisTheBaron/strava-ical/utils/ical"
	"github.com/jaytaylor/html2text"
	"github.com/strava/go.strava"
	"net/http"
	"strings"
	"text/template"
	"time"
)

// ControllerInterface is the interface to which controllers adhere.
type ControllerInterface interface {
	Get()     //method = GET processing
	Post()    //method = POST processing
	Delete()  //method = DELETE processing
	Put()     //method = PUT handling
	Head()    //method = HEAD processing
	Patch()   //method = PATCH treatment
	Options() //method = OPTIONS processing
}

// Controller is the base type of controllers in the strava-ical architecture.
type Controller struct {
	client *strava.Client
	config *structs.Config
}

// Get handles a HTTP GET request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Post handles a HTTP POST request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Delete handles a HTTP DELETE request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Put handles a HTTP PUT request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Put(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Head handles a HTTP HEAD request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Head(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Patch handles a HTTP PATCH request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Patch(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Options handles a HTTP OPTIONS request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Options(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Propfind handles a HTTP PROPFIND request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Propfind(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

func (c *Controller) renderICAL(cal ical.VCalendar, activities []*strava.ActivitySummary, w http.ResponseWriter) {

	nt := template.New("name template")
	nt.Funcs(template.FuncMap{
		"html2text": html2text.FromString,
		"trim":      strings.TrimSpace,
	})
	nt, _ = nt.Parse(c.config.CalendarEntryName)

	dt := template.New("description template")
	dt.Funcs(template.FuncMap{
		"html2text": html2text.FromString,
		"trim":      strings.TrimSpace,
	})
	dt, _ = dt.Parse(c.config.CalendarEntryDescription)

	for _, activity := range activities {

		var name bytes.Buffer
		var desc bytes.Buffer

		err := nt.Execute(&name, struct {
			strava.ActivitySummary
			DistanceKm string
		}{
			DistanceKm: fmt.Sprintf("%.2f", activity.Distance / 1000),
		})

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = dt.Execute(&desc, struct {
			strava.ActivitySummary
			DistanceKm string
		}{
			DistanceKm: fmt.Sprintf("%.2f", activity.Distance / 1000),
		})

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		cal.AddComponent(ical.VComponent(ical.VEvent{
			UID:         fmt.Sprintf("%d", activity.Id),
			SUMMARY:     name.String(),
			DESCRIPTION: desc.String(),
			DTSTART:     activity.StartDateLocal,
			DTEND:       activity.StartDateLocal.Add(time.Duration(activity.ElapsedTime) * time.Second),
			DTSTAMP:     activity.StartDateLocal,
			LOCATION:    activity.City,
			TZID:        "Europe/London",
			AllDay:      false,
		}))

	}

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "inline; filename=strava.ics")
	cal.Encode(w)

}
