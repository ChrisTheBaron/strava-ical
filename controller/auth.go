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

func (a *Auth) OAuthHandler(w http.ResponseWriter, r *http.Request) {
	// you should make this a template in your real application
	fmt.Fprintf(w, `<a href="%s">`, a.authenticator.AuthorizationURL("state1", strava.Permissions.ViewPrivate, true))
	fmt.Fprint(w, `<img src="https://strava.github.io/api/images/btn_connectWith.png" />`)
	fmt.Fprint(w, `</a>`)
}

// OAuthSuccess stores/updates the authenticated user, generates a JWT and stores it in a cookie.
// Then redirects to /calendars.
func (a *Auth) OAuthSuccess(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {

	glog.Infoln("Authenticated successfully")
	glog.Infof("Auth token: %s", auth.AccessToken)

	u := entities.NewUser(auth.Athlete.FirstName, auth.Athlete.LastName, auth.Athlete.Email, auth.Athlete.Id, auth.AccessToken)

	err := a.userModel.Upsert(u)

	if err != nil {
		glog.Error(err)
		http.Redirect(w, r, a.config.Slugs.Login, http.StatusInternalServerError)
		return
	}

	glog.Infoln("Upserted user")

	token, err := utils.GenerateJWT(a.config, u)

	if err != nil {
		glog.Error(err)
		http.Redirect(w, r, a.config.Slugs.Login, http.StatusInternalServerError)
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
		Domain:   a.config.Server.Address,
	})

	// I would use the normal http.Redirect here, but it puts a link for GET requests,
	// which is ugly.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Location", a.config.Slugs.Dashboard)
	w.Write([]byte(fmt.Sprintf("<html><body>Success. Redirecting...<script>window.location = '%s'</script></body></html>", a.config.Slugs.Dashboard)))

	glog.Infoln("Inserted cookie, hopefully.")

}

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
