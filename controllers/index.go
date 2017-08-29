package controllers

import (
	"bytes"
	"github.com/ChrisTheBaron/strava-ical/models"
	"github.com/ChrisTheBaron/strava-ical/structs"
	"github.com/jaytaylor/html2text"
	"github.com/strava/go.strava"
	"net/http"
	"strings"
	"text/template"
)

// IndexController is the controller for the index page.
type IndexController struct {
	Controller
}

// NewIndexController returns a new IndexController with the MyRadio session s
// and configuration context c.
func NewIndexController(cli *strava.Client, con *structs.Config) *IndexController {
	return &IndexController{Controller{client: cli, config: con}}
}

// Get handles the HTTP GET request r for the index page, writing to w.
func (ic *IndexController) Get(w http.ResponseWriter, r *http.Request) {

	im := models.NewIndexModel(ic.client, ic.config)

	activities, err := im.Get()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	cal := ic.config.Calendar

	t := template.New("calendar template")
	t.Funcs(template.FuncMap{
		"html2text": html2text.FromString,
		"trim":      strings.TrimSpace,
	})
	t, _ = t.Parse(ic.config.CalendarDescription)

	var desc bytes.Buffer

	err = t.Execute(&desc, ic.config)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	cal.DESCRIPTION = desc.String()
	cal.X_WR_CALDESC = desc.String()

	ic.renderICAL(cal, activities, w)

}
