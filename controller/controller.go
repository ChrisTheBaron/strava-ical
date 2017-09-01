package controller

import (
	"errors"
	"fmt"
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

// redirect accepts a slug ("login","calendar", etc.) and redirects to an absolute url
// using the rootUrl and protocol in config.
func (c *controller) redirect(w http.ResponseWriter, url string) {
	c.redirectWithStatus(w, url, http.StatusTemporaryRedirect)
}

func (c *controller) redirectWithStatus(w http.ResponseWriter, url string, status int) {
	var fullUrl string
	if url == "/" || url == c.config.RootUrl {
		fullUrl = fmt.Sprintf("%s://%s/", c.config.Protocol, c.config.RootUrl)
	} else {
		fullUrl = fmt.Sprintf("%s://%s/%s", c.config.Protocol, c.config.RootUrl, url)
	}
	w.Header().Set("Location", fullUrl)
	w.WriteHeader(status)
}
