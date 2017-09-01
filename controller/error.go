package controller

import (
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/ChrisTheBaron/strava-ical/utils"
	"github.com/golang/glog"
	"net/http"
)

type Error struct {
	controller
	tr *utils.Template
}

// NewError returns a new controller.Error
func NewError(con *entities.Config, tr *utils.Template) *Error {
	return &Error{controller{config: con}, tr}
}

func (e *Error) E404(w http.ResponseWriter, r *http.Request) {

	if r.Referer() != "" {
		glog.Warningf("404 redirected from %s", r.Referer())
	}

	w.WriteHeader(http.StatusNotFound)
	err := e.tr.Render(w, nil, "errors/404.tmpl")

	if err != nil {
		glog.Error(err)
	}

}

func (e *Error) E500(w http.ResponseWriter, r *http.Request) {

	if r.Referer() != "" {
		glog.Warningf("500 redirected from %s", r.Referer())
	}

	w.WriteHeader(http.StatusInternalServerError)
	err := e.tr.Render(w, nil, "errors/500.tmpl")

	if err != nil {
		glog.Error(err)
	}

}
