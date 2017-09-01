package web

import (
	"database/sql"
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/controller"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/ChrisTheBaron/strava-ical/middleware"
	"github.com/ChrisTheBaron/strava-ical/model"
	"github.com/ChrisTheBaron/strava-ical/services"
	"github.com/ChrisTheBaron/strava-ical/utils"
	"github.com/codegangsta/negroni"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/strava/go.strava"
	"net/http"
	"strings"
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

	/*
	   ------------------------------------------
	   				UTILS
	   ------------------------------------------
	*/

	tr := utils.NewTemplate(c)

	/*
		------------------------------------------
		               FACTORIES
		------------------------------------------
	*/

	sf := services.NewStravaFactory(c)

	/*
		------------------------------------------
		                MODELS
		------------------------------------------
	*/

	um := model.NewUser(c, db)
	cm := model.NewCalendar(c, db)

	/*
	   ------------------------------------------
	                 MIDDLEWARE
	   ------------------------------------------
	*/

	am := middleware.NewVerifyJWT(um, db, c)

	authenticator := strava.OAuthAuthenticator{
		CallbackURL:            fmt.Sprintf("http://%s/%s", c.RootUrl, c.Slugs.OAuthCallback),
		RequestClientGenerator: nil,
	}

	autoClbPath, err := authenticator.CallbackPath()

	if err != nil {
		return nil, err
	}

	/*
		------------------------------------------
		                CONTROLLERS
		------------------------------------------
	*/

	ec := controller.NewError(c, tr)
	cc := controller.NewCalendar(c, cm, um, sf, tr)
	ac := controller.NewAuth(c, um, authenticator)
	ic := controller.NewIndex(c, tr)

	/*
		------------------------------------------
		                ROUTES
		------------------------------------------
	*/

	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(ec.E404)

	getRouter := router.Methods("GET").Subrouter()
	postRouter := router.Methods("POST").Subrouter()
	//deleteRouter := router.Methods("DELETE").Subrouter()

	getRouter.Handle("/", http.HandlerFunc(ic.Get))

	getRouter.PathPrefix("/static").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := strings.TrimLeft(r.URL.Path, "/")
		glog.Infoln(url)
		if asset, err := utils.Asset(url); err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(asset)
		} else {
			glog.Error(err)
			w.WriteHeader(http.StatusNotFound)
		}
	})

	getRouter.Handle(fmt.Sprintf("/%s", c.Slugs.OAuth), http.HandlerFunc(ac.OAuthHandler))
	getRouter.Handle(autoClbPath, authenticator.HandlerFunc(ac.OAuthSuccess, ac.OAuthFailure))

	// /calendar/
	// list all
	getRouter.HandleFunc(fmt.Sprintf("/%s", c.Slugs.Calendars), am(http.HandlerFunc(cc.Get)))

	// create
	postRouter.HandleFunc(fmt.Sprintf("/%s", c.Slugs.Calendars), am(http.HandlerFunc(cc.Post)))

	// list
	getRouter.HandleFunc(fmt.Sprintf("/%s/{id:.{36}}.ics", c.Slugs.Calendars), am(http.HandlerFunc(cc.GetICALById)))
	getRouter.HandleFunc(fmt.Sprintf("/%s/{id:.{36}}", c.Slugs.Calendars), am(http.HandlerFunc(cc.GetById)))

	s.UseHandler(router)

	return &s, nil

}

func connectToDB(config *entities.Config) (*sql.DB, error) {
	return sql.Open("sqlite3", config.DBPath)
}
