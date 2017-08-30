package web

import (
	"database/sql"
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/controller"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/ChrisTheBaron/strava-ical/middleware"
	"github.com/ChrisTheBaron/strava-ical/model"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/strava/go.strava"
	"net/http"
)

type Server struct {
	*negroni.Negroni
}

func NewServer(c *entities.Config) (*Server, error) {

	strava.ClientId = c.StravaClientId
	strava.ClientSecret = c.StravaClientSecret

	db, err := connectToDB(c)

	if err != nil {
		return nil, err
	}

	s := Server{negroni.Classic()}

	um := model.NewUser(c, db)

	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(http.NotFound)

	getRouter := router.Methods("GET").Subrouter()

	// Routes go in here

	cc := controller.NewCalendar(c)

	getRouter.HandleFunc(c.Slugs.Dashboard, middleware.VerifyJWT(um, db, c, http.HandlerFunc(cc.Get)))

	//ic := controllers.NewIndexController(client, c)
	//getRouter.HandleFunc("/strava.ics", ic.Get)

	authenticator := strava.OAuthAuthenticator{
		CallbackURL:            fmt.Sprintf("http://%s%s", c.Server.Address, c.Slugs.OAuthCallback),
		RequestClientGenerator: nil,
	}

	ac := controller.NewAuth(c, um, authenticator)

	path, err := authenticator.CallbackPath()

	if err != nil {
		return nil, err
	}

	getRouter.Handle(c.Slugs.OAuth, http.HandlerFunc(ac.OAuthHandler))
	getRouter.Handle(path, authenticator.HandlerFunc(ac.OAuthSuccess, ac.OAuthFailure))

	// End routes

	s.UseHandler(router)

	return &s, nil

}

func connectToDB(config *entities.Config) (*sql.DB, error) {
	return sql.Open("sqlite3", config.DBPath)
}
