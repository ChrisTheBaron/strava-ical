package controller

import (
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/ChrisTheBaron/strava-ical/model"
	"github.com/ChrisTheBaron/strava-ical/utils"
	"github.com/golang/glog"
	"github.com/strava/go.strava"
	"net/http"
	"time"
)

// Auth is the controller for the OAuth pages.
type Auth struct {
	controller
	userModel     *model.User
	authenticator strava.OAuthAuthenticator
}

// NewAuth returns a new controller.Auth
func NewAuth(con *entities.Config, um *model.User, authenticator strava.OAuthAuthenticator) *Auth {
	return &Auth{controller{config: con}, um, authenticator}
}

// Login checks for a cookie/header, otherwise redirects to Strava for auth
func (a *Auth) OAuthHandler(w http.ResponseWriter, r *http.Request) {
	if token := utils.GetTokenFromRequest(a.config, r); token != "" {
		a.redirect(w, a.config.Slugs.Calendars)
	} else {
		http.Redirect(w, r,
			a.authenticator.AuthorizationURL("state1", strava.Permissions.ViewPrivate, true), http.StatusTemporaryRedirect)
	}
}

// OAuthSuccess stores/updates the authenticated user, generates c JWT and stores it in c cookie.
// Then redirects to /calendars.
func (a *Auth) OAuthSuccess(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {

	glog.Infoln("Authenticated successfully")
	glog.Infof("Auth token: %s", auth.AccessToken)

	u := entities.User{
		Firstname:         auth.Athlete.FirstName,
		Lastname:          auth.Athlete.LastName,
		StravaId:          auth.Athlete.Id,
		Email:             auth.Athlete.Email,
		StravaAccessToken: auth.AccessToken,
	}

	err := a.userModel.Upsert(u)

	if err != nil {
		glog.Error(err)
		a.redirect(w, "/")
		return
	}

	glog.Infoln("Upserted user")

	token, err := utils.GenerateJWT(a.config, u)

	if err != nil {
		glog.Error(err)
		a.redirect(w, "/")
		return
	}

	glog.Infof("Generated JWT: %s", token)

	expiration := time.Now().Add(365 * 24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Value:    token,
		Name:     a.config.JWTCookieName,
		HttpOnly: true,
		Secure:   false,
		Expires:  expiration,
		Path:     "/",
		Domain:   a.config.RootUrl,
	})

	a.redirect(w, a.config.Slugs.Calendars)

	glog.Infoln("Inserted cookie, hopefully.")

}

// OAuthFailure redirects to login/ with an error messsage
// @TODO: Do that
func (a *Auth) OAuthFailure(err error, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authorization Failure:")

	// some standard error checking
	if err == strava.OAuthAuthorizationDeniedErr {
		fmt.Fprintln(w, "The user clicked the 'Do not Authorize' button on the previous page.")
		fmt.Fprintln(w, "This is the main error your application should handle.")
	} else if err == strava.OAuthInvalidCredentialsErr {
		fmt.Fprintln(w, "You provided an incorrect client_id or client_secret.\nDid you remember to set them at the begininng of this file?")
	} else if err == strava.OAuthInvalidCodeErr {
		fmt.Fprintln(w, "The temporary token was not recognized, this shouldn't happen normally")
	} else if err == strava.OAuthServerErr {
		fmt.Fprintln(w, "There was some sort of server error, try again to see if the problem continues")
	} else {
		fmt.Fprintln(w, err)
	}
}
