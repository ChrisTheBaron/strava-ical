package web

import (
	"github.com/ChrisTheBaron/strava-ical/controllers"
	"github.com/ChrisTheBaron/strava-ical/structs"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/strava/go.strava"
	"net/http"
)

type Server struct {
	*negroni.Negroni
}

func NewServer(c *structs.Config) (*Server, error) {

	s := Server{negroni.Classic()}

	client := strava.NewClient(c.StravaApiKey)

	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(http.NotFound)

	getRouter := router.Methods("GET").Subrouter()

	// Routes go in here

	ic := controllers.NewIndexController(client, c)
	getRouter.HandleFunc("/strava.ics", ic.Get)

	// End routes

	s.UseHandler(router)

	return &s, nil

}
